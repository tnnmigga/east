package module

import (
	"fmt"
	"reflect"
)

type MessageHandler func(msg any)

func RegisterHandler[T any](m *Module, msg any, cb func(msg T)) {
	msgType := reflect.TypeOf(msg)
	_, ok := m.handlers[msgType]
	if ok {
		panic(fmt.Errorf("RegMsgCb duplicated reg %v", msgType))
	}
	m.handlers[msgType] = func(msg0 any) {
		msg := msg0.(T)
		cb(msg)
	}
}
