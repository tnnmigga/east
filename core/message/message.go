package message

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
	buffer  chan any
	recvers map[reflect.Type]IRecver
)

func init() {
	once.Do(func() {
		recvers = map[reflect.Type]IRecver{}
	})
}

type IRecver interface {
	Name() string
	MQ() chan any
}

// Init
func Init() {
	once.Do(func() {
		buffer = make(chan any, iconf.Int32("mq-len", 100000))
		go util.ExecAndRecover(dispatch)
	})
}

func Cast(serverID uint32, msg any, byStream ...bool) {
	useStream := util.FirstOrDefault(byStream, true) // 默认使用流
	if serverID == iconf.ServerID() {
		recv, ok := recvers[reflect.TypeOf(msg)]
		if !ok {
			log.Errorf("message cast recv not fuound %v", util.ReflectName(msg))
			return
		}
		select {
		case recv.MQ() <- msg:
		case buffer <- msg: // 一次重试机会
			log.Errorf("metssage cast mq full %s %s", recv.Name(), util.ReflectName(msg))
		default:
			log.Errorf("message cast faild, buffer full %s %s", util.ReflectName(msg), util.String(msg))
		}
		return
	}
	if useStream {
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
	recvers[mType] = recver
}

func dispatch() {
	for msg := range buffer {
		recver, ok := recvers[reflect.TypeOf(msg)]
		if !ok {
			log.Errorf("message dispatch recver not found %v", util.ReflectName(msg))
			return
		}
		select {
		case recver.MQ() <- msg:
		default:
			log.Errorf("message dispatch faild, mq full %s %s %s", recver.Name(), util.ReflectName(msg), util.String(msg))
		}
	}
}
