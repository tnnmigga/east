package message

import (
	"east/core/define"
	"east/core/iconf"
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

type IRecver interface {
	Name() string
	MQ() chan any
}

// Attach 放置全局分发入口
func EnableBuffer() {
	once.Do(func() {
		buffer = make(chan any, iconf.Int32("mq-len", 100000))
		go util.ExecAndRecover(dispatch)
	})
}

func Cast(serverID uint32, msg any) {
	if serverID != iconf.ServerID() {
		Cast(iconf.ServerID(), &define.Package{
			ServerID: serverID,
			Body:     msg,
		})
		return
	}
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
}

func Broadcast(serverType string, msg any) {
	pkg := &define.BroadcastPackage{
		ServerType: serverType,
		Body:       msg,
	}
	Cast(iconf.ServerID(), pkg)
}

func RPC[T any](module define.IModule, serverID uint32, req any, cb func(resp T, err error)) {
	pkg := &define.RPCRequest{
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

func RegisterRecver[msg any](recver IRecver) {
	mType := reflect.TypeOf(new(msg))
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
