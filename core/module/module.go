package module

import (
	"east/core/idef"
	"east/core/log"
	"east/core/util"
	"fmt"
	"reflect"
	"runtime/debug"
)

const (
	DefaultMQLen = 100000
)

var (
	rpcPackage = reflect.TypeOf((*idef.RPCPackage)(nil))
	rpcRequest = reflect.TypeOf((*idef.RPCRequest)(nil))
)

type Module struct {
	name      string
	mq        chan any
	handlers  map[reflect.Type]*idef.Handler
	hooks     [idef.ServerStateMax][2][]func() error
	closeSign chan struct{}
}

func New(name string, mqLen int32) *Module {
	return &Module{
		name:      name,
		mq:        make(chan any, mqLen),
		handlers:  map[reflect.Type]*idef.Handler{},
		closeSign: make(chan struct{}, 1),
	}
}

func (m *Module) Name() string {
	return m.name
}

func (m *Module) MQ() chan any {
	return m.mq
}

func (m *Module) RegisterHandler(mType reflect.Type, handler *idef.Handler) {
	_, ok := m.handlers[mType]
	if ok {
		// 一个module内一个msg只能被注册一次, 但不同模块可以分别注册监听同一个消息
		panic(fmt.Errorf("RegisterHandler multiple registration %v", mType))
	}
	m.handlers[mType] = handler
}

func (m *Module) Hook(state idef.ServerState, stage int) []func() error {
	return m.hooks[state][stage]
}

func (m *Module) Before(state idef.ServerState, hook func() error) {
	m.hooks[state][0] = append(m.hooks[state][0], hook)
}

func (m *Module) After(state idef.ServerState, hook func() error) {
	if state >= idef.ServerStateClose {
		panic("module after close hook not support")
	}
	m.hooks[state][1] = append(m.hooks[state][1], hook)
}

func (m *Module) Run() {
	defer func() {
		log.Infof("%v has stoped", m.Name())
		m.closeSign <- struct{}{}
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
	<-m.closeSign
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
