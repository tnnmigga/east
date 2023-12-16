package msgbus

import (
	"east/core/codec"
	"east/core/idef"
	"fmt"
	"reflect"
)

func RegisterHandler[T any](m idef.IModule, fn func(msg T)) {
	mValue := *new(T)
	mType := reflect.TypeOf(mValue)
	codec.Register(mValue)
	registerRecver(mType, m)
	m.RegisterHandler(mType, &idef.Handler{
		Cb: func(msg0 any) {
			msg := msg0.(T)
			fn(msg)
		},
	})
}

func RegisterRPC[T any](m idef.IModule, fn func(msg T, resp func(msg any))) {
	msgType := reflect.TypeOf(new(T))
	m.RegisterHandler(msgType, &idef.Handler{
		RPC: func(msg0 any, resp func(any)) {
			msg := msg0.(T)
			fn(msg, resp)
		},
	})
}

func registerRecver(mType reflect.Type, recver IRecver) {
	if _, has := recvers[mType]; has {
		panic(fmt.Errorf("message has registered %v", mType.Elem().Name()))
	}
	recvers[mType] = append(recvers[mType], recver)
}
