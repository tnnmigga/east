package message

import (
	"east/core/define"
	"fmt"
	"reflect"
)

func RegisterHandler[T any](m define.IModule, fn func(msg T)) {
	RegisterRecver[T](m)
	msgType := reflect.TypeOf(new(T))
	_, ok := m.Handlers()[msgType]
	if ok {
		panic(fmt.Errorf("RegisterHandler multiple registration %v", msgType))
	}
	m.Handlers()[msgType] = &define.HandlerFn{
		Cb: func(msg0 any) {
			msg := msg0.(T)
			fn(msg)
		},
	}
}

func RegisterRPC[T any](m define.IModule, fn func(msg T, resp func(msg any))) {
	msgType := reflect.TypeOf(new(T))
	_, ok := m.Handlers()[msgType]
	if ok {
		panic(fmt.Errorf("RegisterHandler multiple registration %v", msgType))
	}
	m.Handlers()[msgType] = &define.HandlerFn{
		RPC: func(msg0 any, resp func(any)) {
			msg := msg0.(T)
			fn(msg, resp)
		},
	}
}
