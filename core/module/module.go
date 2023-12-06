package module

import (
	"east/core/idef"
	"east/core/log"
	"east/core/util"
	"fmt"
	"reflect"
	"runtime/debug"
)

type Module struct {
	name     string
	mq       chan any
	handlers map[reflect.Type]*idef.HandlerFn
	closeSig chan struct{}
}

func New(name string, mqLen int32) *Module {
	return &Module{
		name:     name,
		mq:       make(chan any, mqLen),
		handlers: map[reflect.Type]*idef.HandlerFn{},
		closeSig: make(chan struct{}, 1),
	}
}

func (m *Module) Name() string {
	return m.name
}

func (m *Module) MQ() chan any {
	return m.mq
}

func (m *Module) Handlers() map[reflect.Type]*idef.HandlerFn {
	return m.handlers
}

func (m *Module) Run() {
	defer util.RecoverPanic()
	defer func() {
		log.Infof("%v has stoped", m.Name())
		m.closeSig <- struct{}{}
	}()
	for msg := range m.MQ() {
		msgType := reflect.TypeOf(msg)
		switch msgType {
		case rpcPackage: // 被发起rpc
			m.rpc(msg.(*idef.RPCPackage))
		case rpcRequest: // rpc请求完成
			m.rpcResp(msg.(*idef.RPCRequest))
		default:
			m.cb(msg)
		}
	}
}

func (m *Module) Stop() {
	log.Infof("try stop %s", m.name)
	close(m.mq)
	<-m.closeSig
	// log.Infof("stop %s success", m.name)
}

func (m *Module) cb(msg any) {
	defer util.RecoverPanic()
	msgType := reflect.TypeOf(msg)
	fns, ok := m.handlers[msgType]
	if !ok {
		log.Errorf("handler not exist %v", msgType)
		return
	}
	fns.Cb(msg)
}

func (m *Module) rpc(msg *idef.RPCPackage) {
	defer func() {
		if r := recover(); r != nil {
			msg.Err <- fmt.Errorf("%v: %s", r, debug.Stack())
		}
	}()
	msgType := reflect.TypeOf(msg.Req)
	fns, ok := m.handlers[msgType]
	if !ok {
		msg.Err <- fmt.Errorf("rpc handler not found %v", msgType)
		return
	}
	fns.RPC(msg.Req, func(v any) {
		msg.Resp <- v
	})
}

func (m *Module) rpcResp(req *idef.RPCRequest) {
	defer util.RecoverPanic()
	req.Cb(req.Resp, req.Err)
}
