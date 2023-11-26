package module

import (
	"eden/core/message"
	"fmt"
	"reflect"
)

type HandlerFn struct {
	Cb  func(msg any)
	RPC func(msg any, resp func(any))
}

var (
	rpcPackage = reflect.TypeOf((*message.RPCPackage)(nil))
	rpcResp    = reflect.TypeOf((*message.RPCRequest)(nil))
)

func RegisterHandler[T any](m *Module, fn func(msg T)) {
	msgType := reflect.TypeOf(new(T))
	_, ok := m.handlers[msgType]
	if ok {
		panic(fmt.Errorf("RegisterHandler multiple registration %v", msgType))
	}
	m.handlers[msgType] = &HandlerFn{
		Cb: func(msg0 any) {
			msg := msg0.(T)
			fn(msg)
		},
	}
}

func RegisterRPC[T any](m *Module, fn func(msg T, respFn func(resp any))) {
	msgType := reflect.TypeOf(new(T))
	_, ok := m.handlers[msgType]
	if ok {
		panic(fmt.Errorf("RegisterHandler multiple registration %v", msgType))
	}
	m.handlers[msgType] = &HandlerFn{
		RPC: func(msg0 any, resp func(any)) {
			msg := msg0.(T)
			fn(msg, resp)
		},
	}
}
