package msgbus

import (
	"east/core/iconf"
	"east/core/idef"
	"east/core/log"
	"east/core/util"

	"reflect"
)

var (
	recvers map[reflect.Type][]IRecver
)

func init() {
	recvers = map[reflect.Type][]IRecver{}
}

type IRecver interface {
	Name() string
	MQ() chan any
}

func Cast(serverID uint32, msg any, opts ...castOpt) {
	if serverID == iconf.ServerID() {
		messageDispatch(msg, opts...)
		return
	}
	if nonuse, find := findCastOpt[bool](opts, keyNonuseStream); nonuse && find { // 不使用流
		Cast(iconf.ServerID(), &idef.CastPackage{
			ServerID: serverID,
			Body:     msg,
		}, opts...)
		return
	}
	Cast(iconf.ServerID(), &idef.StreamCastPackage{
		ServerID: serverID,
		Body:     msg,
	}, opts...)
}

func Broadcast(serverType string, msg any) {
	pkg := &idef.BroadcastPackage{
		ServerType: serverType,
		Body:       msg,
	}
	Cast(iconf.ServerID(), pkg)
}

func RPC[T any](module idef.IModule, serverID uint32, req any, cb func(resp T, err error)) {
	pkg := &idef.RPCRequest{
		Module:   module,
		ServerID: serverID,
		Req:      req,
		Cb:       warpRPCCb(cb),
	}
	Cast(iconf.ServerID(), pkg)
}

func warpRPCCb[T any](cb func(resp T, err error)) func(resp any, err error) {
	return func(pkg any, err error) {
		resp := pkg.(T)
		cb(resp, err)
	}
}

func messageDispatch(msg any, opts ...castOpt) {
	recvs, ok := recvers[reflect.TypeOf(msg)]
	modName, oneOfMod := findCastOpt[string](opts, keyOneOfModules)
	var castSucc bool
	for _, recv := range recvs {
		if !ok {
			log.Errorf("message cast recv not fuound %v", util.StructName(msg))
			return
		}
		if oneOfMod && modName != recv.Name() {
			continue
		}
		select {
		case recv.MQ() <- msg:
			castSucc = true
		default:
			onCastFail(recv, msg)
		}
	}
	if castSucc {
		return
	}
	log.Errorf("metssage cast mq full %s %s", util.StructName(msg), util.String(msg))
}

// onCastFail 消息投递失败处理
func onCastFail(recver IRecver, msg any) {
	log.Errorf("message cast faild, %s %s", util.StructName(msg), util.String(msg))
}
