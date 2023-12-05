package message

import (
	"east/core/idef"
	"fmt"
	"reflect"
)

func RegisterHandler[T any](m idef.IModule, fn func(msg T)) {
	mType := reflect.TypeOf(new(T))
	registerRecver(mType, m)
	_, ok := m.Handlers()[mType]
	if ok {
		panic(fmt.Errorf("RegisterHandler multiple registration %v", mType))
	}
	m.Handlers()[mType] = &idef.HandlerFn{
		Cb: func(msg0 any) {
			msg := msg0.(T)
			fn(msg)
		},
	}
}

func RegisterRPC[T any](m idef.IModule, fn func(msg T, resp func(msg any))) {
	msgType := reflect.TypeOf(new(T))
	_, ok := m.Handlers()[msgType]
	if ok {
		panic(fmt.Errorf("RegisterHandler multiple registration %v", msgType))
	}
	m.Handlers()[msgType] = &idef.HandlerFn{
		RPC: func(msg0 any, resp func(any)) {
			msg := msg0.(T)
			fn(msg, resp)
		},
	}
}
