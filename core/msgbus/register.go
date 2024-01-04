package msgbus

import (
	"east/core/codec"
	"east/core/idef"
	"fmt"
	"log"
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

func RegisterRPC[T any](m idef.IModule, fn func(msg T, resp func(any), err func(error))) {
	mValue := *new(T)
	mType := reflect.TypeOf(mValue)
	codec.Register(mValue)
	registerRecver(mType, m)
	m.RegisterHandler(mType, &idef.Handler{
		RPC: func(msg0 any, r func(any), e func(error)) {
			msg := msg0.(T)
			fn(msg, r, e)
		},
	})
}

func registerRecver(mType reflect.Type, recver IRecver) {
	if _, has := recvers[mType]; has {
		log.Fatal(fmt.Errorf("message has registered %v", mType.Elem().Name()))
	}
	recvers[mType] = append(recvers[mType], recver)
}
