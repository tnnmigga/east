package msgbus

import (
	"east/core/codec"
	"east/core/idef"
	"fmt"
	"log"
	"reflect"
)

func RegisterHandler[T any](com idef.IComponent, fn func(T)) {
	mValue := *new(T)
	mType := reflect.TypeOf(mValue)
	codec.Register(mValue)
	registerRecver(mType, com)
	com.RegisterHandler(mType, &idef.Handler{
		Cb: func(msg0 any) {
			msg := msg0.(T)
			fn(msg)
		},
	})
}

func RegisterRPC[T any](com idef.IComponent, fn func(msg T, resolve func(any), reject func(error))) {
	mValue := *new(T)
	mType := reflect.TypeOf(mValue)
	codec.Register(mValue)
	registerRecver(mType, com)
	com.RegisterHandler(mType, &idef.Handler{
		RPC: func(msg0 any, res func(any), rej func(error)) {
			msg := msg0.(T)
			fn(msg, res, rej)
		},
	})
}

func registerRecver(mType reflect.Type, recver IRecver) {
	if _, has := recvers[mType]; has {
		log.Fatal(fmt.Errorf("message has registered %v", mType.Elem().Name()))
	}
	recvers[mType] = append(recvers[mType], recver)
}
