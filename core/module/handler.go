package module

import (
	"fmt"
	"reflect"
)

type MessageHandler func(msg any)

func RegisterHandler[T any](m IModule, cb func(msg T)) {
	msgType := reflect.TypeOf(new(T))
	handlers := m.Handlers()
	_, ok := handlers[msgType]
	if ok {
		panic(fmt.Errorf("RegMsgCb duplicated reg %v", msgType))
	}
	handlers[msgType] = func(msg0 any) {
		msg := msg0.(T)
		cb(msg)
	}
}
