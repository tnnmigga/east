package compt

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
	rpcReqType  = reflect.TypeOf((*idef.RPCRequest)(nil))
	rpcRespType = reflect.TypeOf((*idef.RPCResponse)(nil))
)

type Component struct {
	name      string
	mq        chan any
	handlers  map[reflect.Type]*idef.Handler
	hooks     [idef.ServerStateClose + 1][2][]func() error
	closeSign chan struct{}
}

func New(name string, mqLen int32) *Component {
	return &Component{
		name:      name,
		mq:        make(chan any, mqLen),
		handlers:  map[reflect.Type]*idef.Handler{},
		closeSign: make(chan struct{}, 1),
	}
}

func (com *Component) Name() string {
	return com.name
}

func (com *Component) MQ() chan any {
	return com.mq
}

func (com *Component) Assign(msg any) {
	select {
	case com.mq <- msg:
	default:
		log.Errorf("modele %s mq full, lose %s", com.name, util.String(msg))
	}
}

func (com *Component) RegisterHandler(mType reflect.Type, handler *idef.Handler) {
	_, ok := com.handlers[mType]
	if ok {
		// 一个module内一个msg只能被注册一次, 但不同模块可以分别注册监听同一个消息
		log.Fatal(fmt.Errorf("RegisterHandler multiple registration %v", mType))
	}
	com.handlers[mType] = handler
}

func (com *Component) Hook(state idef.ServerState, stage int) []func() error {
	return com.hooks[state][stage]
}

func (com *Component) Before(state idef.ServerState, hook func() error) {
	if state <= idef.ServerStateInit {
		panic("component after close hook not support")
	}
	com.hooks[state][0] = append(com.hooks[state][0], hook)
}

func (com *Component) After(state idef.ServerState, hook func() error) {
	if state >= idef.ServerStateClose {
		panic("component after close hook not support")
	}
	com.hooks[state][1] = append(com.hooks[state][1], hook)
}

func (com *Component) Run() {
	defer func() {
		log.Infof("%v has stoped", com.Name())
		com.closeSign <- struct{}{}
	}()
	for msg := range com.mq {
		msgType := reflect.TypeOf(msg)
		switch msgType {
		case rpcReqType: // 被发起rpc
			com.rpc(msg.(*idef.RPCRequest))
		case rpcRespType: // rpc请求完成
			com.rpcResp(msg.(*idef.RPCResponse))
		default:
			com.cb(msg)
		}
	}
}

func (com *Component) Stop() {
	log.Infof("try stop %s", com.name)
	close(com.mq)
	<-com.closeSign
}

func (com *Component) cb(msg any) {
	defer util.RecoverPanic()
	msgType := reflect.TypeOf(msg)
	h, ok := com.handlers[msgType]
	if !ok {
		log.Errorf("handler not exist %v", msgType)
		return
	}
	h.Cb(msg)
}

func (com *Component) rpc(msg *idef.RPCRequest) {
	defer func() {
		if r := recover(); r != nil {
			msg.Err <- fmt.Errorf("%v: %s", r, debug.Stack())
		}
	}()
	msgType := reflect.TypeOf(msg.Req)
	h, ok := com.handlers[msgType]
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

func (com *Component) rpcResp(req *idef.RPCResponse) {
	defer util.RecoverPanic()
	req.Cb(req.Resp, req.Err)
}
