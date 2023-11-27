package module

import (
	"eden/core/log"
	"eden/core/message"
	"eden/core/util"
	"fmt"
	"reflect"
	"runtime/debug"
)

type IModule interface {
	Name() string
	MQ() chan any
	Init() error
	Run()
	Close()
}

type Module struct {
	name     string
	mq       chan any
	handlers map[reflect.Type]*HandlerFn
	closeSig chan struct{}
}

var modules = map[int64]*Module{}

func NewModule(mType string, mqLen int32) *Module {
	m := &Module{
		name:     mType,
		mq:       make(chan any, mqLen),
		handlers: map[reflect.Type]*HandlerFn{},
		closeSig: make(chan struct{}, 1),
	}
	return m
}

func (m *Module) Name() string {
	return m.name
}

func (m *Module) MQ() chan any {
	return m.mq
}

func (m *Module) Init() error {
	return nil
}

func (m *Module) Run() {
	defer util.RecoverPanic()
	defer func() {
		log.Infof("module %v has stoped", m.Name())
		m.closeSig <- struct{}{}
	}()
	for msg := range m.MQ() {
		msgType := reflect.TypeOf(msg)
		switch msgType {
		case rpcPackage: // 被发起rpc
			m.rpc(msg.(*message.RPCPackage))
		case rpcRequest: // rpc请求完成
			m.rpcResp(msg.(*message.RPCRequest))
		default:
			m.cb(msg)
		}
	}
}

func (m *Module) Close() {
	close(m.mq)
	<-m.closeSig
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

func (m *Module) rpc(msg *message.RPCPackage) {
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

func (m *Module) rpcResp(req *message.RPCRequest) {
	defer util.RecoverPanic()
	req.Cb(req.Resp, req.Err)
}
