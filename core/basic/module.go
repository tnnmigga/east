package basic

import (
	"east/core/idef"
	"east/core/msgbus"
	"east/core/util"
	"east/core/zlog"
	"fmt"
	"reflect"
	"runtime/debug"
)

const (
	DefaultMQLen = 100000
)

type Module struct {
	name      string
	mq        chan any
	handlers  map[reflect.Type]*idef.Handler
	hooks     [idef.ServerStateExit + 1][2][]func() error
	closeSign chan struct{}
}

func New(name string, mqLen int32) *Module {
	m := &Module{
		name:      name,
		mq:        make(chan any, mqLen),
		handlers:  map[reflect.Type]*idef.Handler{},
		closeSign: make(chan struct{}, 1),
	}
	msgbus.RegisterHandler(m, m.onRPCRequest)
	msgbus.RegisterHandler(m, m.onRPCResponse)
	return m
}

func (m *Module) Name() string {
	return m.name
}

func (m *Module) MQ() chan any {
	return m.mq
}

func (m *Module) Assign(msg any) {
	select {
	case m.mq <- msg:
	default:
		zlog.Errorf("modele %s mq full, lose %s", m.name, util.String(msg))
	}
}

func (m *Module) RegisterHandler(mType reflect.Type, handler *idef.Handler) {
	_, ok := m.handlers[mType]
	if ok {
		// 一个module内一个msg只能被注册一次, 但不同模块可以分别注册监听同一个消息
		zlog.Fatal(fmt.Errorf("RegisterHandler multiple registration %v", mType))
	}
	m.handlers[mType] = handler
}

func (m *Module) Hook(state idef.ServerState, stage int) []func() error {
	return m.hooks[state][stage]
}

func (m *Module) Before(state idef.ServerState, hook func() error) {
	if state <= idef.ServerStateInit {
		zlog.Fatal("module after close hook not support")
	}
	m.hooks[state][0] = append(m.hooks[state][0], hook)
}

func (m *Module) After(state idef.ServerState, hook func() error) {
	if state >= idef.ServerStateExit {
		zlog.Fatal("module after close hook not support")
	}
	m.hooks[state][1] = append(m.hooks[state][1], hook)
}

func (m *Module) Run() {
	defer func() {
		zlog.Infof("%v has stoped", m.Name())
		m.closeSign <- struct{}{}
	}()
	for msg := range m.mq {
		m.cb(msg)
	}
}

func (m *Module) Stop() {
	zlog.Infof("try stop %s", m.name)
	close(m.mq)
	<-m.closeSign
}

func (m *Module) cb(msg any) {
	defer util.RecoverPanic()
	msgType := reflect.TypeOf(msg)
	h, ok := m.handlers[msgType]
	if !ok {
		zlog.Errorf("handler not exist %v", msgType)
		return
	}
	h.Cb(msg)
}

func (m *Module) onRPCRequest(msg *idef.RPCRequest) {
	defer func() {
		if r := recover(); r != nil {
			msg.Err <- fmt.Errorf("%v: %s", r, debug.Stack())
		}
	}()
	msgType := reflect.TypeOf(msg.Req)
	h, ok := m.handlers[msgType]
	if !ok {
		msg.Err <- fmt.Errorf("rpc handler not found %v", msgType)
		return
	}
	h.RPC(msg.Req, func(v any) {
		msg.Resp <- v
	}, func(err error) {
		msg.Err <- err
	})
}

func (m *Module) onRPCResponse(req *idef.RPCResponse) {
	defer util.RecoverPanic()
	req.Cb(req.Resp, req.Err)
}
