package message

import (
	"eden/core/iconf"
	"eden/core/infra"
	"sync"
)

var (
	once sync.Once
	mq   chan<- *Package
)

// Attach 放置全局分发入口
func Attach(c chan<- *Package) {
	once.Do(func() {
		mq = c
	})
}

func Cast(serverID uint32, module string, msg any) {
	pkg := &Package{
		ServerID: serverID,
		Module:   module,
		Body:     msg,
	}
	if serverID != iconf.ServerID() {
		Cast(iconf.ServerID(), infra.Nats, pkg)
		return
	}
	mq <- pkg
}

func Broadcast(serverType string, module string, msg any) {
	pkg := &BroadcastPackage{
		ServerType: serverType,
		Module:     module,
		Body:       msg,
	}
	Cast(iconf.ServerID(), infra.Nats, pkg)
}

func RPC[T any](caller string, serverID uint32, module string, req any, cb func(resp T, err error)) {
	pkg := &RPCRequest{
		Caller:   caller,
		ServerID: serverID,
		Module:   module,
		Req:      req,
		Cb:       rpcCb(cb),
	}
	Cast(iconf.ServerID(), infra.Nats, pkg)
}

func rpcCb[T any](cb func(resp T, err error)) func(resp any, err error) {
	return func(pkg any, err error) {
		resp := pkg.(T)
		cb(resp, err)
	}
}
