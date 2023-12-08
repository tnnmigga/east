package msgbus

import (
	"east/core/iconf"
	"east/core/idef"
	"east/core/log"
	"east/core/util"

	"fmt"
	"reflect"
	"sync"
)

var (
	once    sync.Once
	recvers map[reflect.Type][]IRecver
)

func init() {
	once.Do(func() {
		recvers = map[reflect.Type][]IRecver{}
	})
}

type IRecver interface {
	Name() string
	MQ() chan any
}

func Cast(serverID uint32, msg any, byStream ...bool) {
	if serverID == iconf.ServerID() {
		messageDispatch(msg)
		return
	}
	if util.FirstOrDefault(byStream, true) { // 默认使用流
		Cast(iconf.ServerID(), &idef.StreamCastPackage{
			ServerID: serverID,
			Body:     msg,
		})
		return
	}
	Cast(iconf.ServerID(), &idef.CastPackage{
		ServerID: serverID,
		Body:     msg,
	})
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

func registerRecver(mType reflect.Type, recver IRecver) {
	if _, has := recvers[mType]; has {
		panic(fmt.Errorf("message has registered %v", mType.Elem().Name()))
	}
	recvers[mType] = append(recvers[mType], recver)
}

func messageDispatch(msg any) {
	recvs, ok := recvers[reflect.TypeOf(msg)]
	for _, recv := range recvs {
		if !ok {
			log.Errorf("message cast recv not fuound %v", util.StructName(msg))
			return
		}
		select {
		case recv.MQ() <- msg:
		default:
			log.Errorf("metssage cast mq full %s %s", recv.Name(), util.StructName(msg))
			onCastFail(recv, msg)
		}
	}
}

// onCastFail 消息投递失败处理
func onCastFail(recver IRecver, msg any) {
	log.Errorf("message cast faild, %s %s", util.StructName(msg), util.String(msg))
}
