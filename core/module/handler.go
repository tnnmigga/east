package module

import (
	"fmt"
	"reflect"

	"github.com/gogo/protobuf/proto"
)

type MessageHandler func(msg proto.Message)

func RegisterHandler[T proto.Message](m *Module, msg proto.Message, cb func(msg T)) {
	msgType := reflect.TypeOf(msg)
	_, ok := m.handlers[msgType]
	if ok {
		panic(fmt.Errorf("RegMsgCb duplicated reg %v", msgType))
	}
	m.handlers[msgType] = func(msg0 proto.Message) {
		msg := msg0.(T)
		cb(msg)
	}
}
