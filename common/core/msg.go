package core

import (
	"fmt"
	"reflect"

	"github.com/gogo/protobuf/proto"
)

type MsgCb func(msg proto.Message)

var msgCbs = map[reflect.Type]MsgCb{}

func RegMsgCb(msg proto.Message, cb MsgCb) {
	msgType := reflect.TypeOf(msg)
	_, ok := msgCbs[msgType]
	if ok {
		panic(fmt.Errorf("RegMsgCb duplicated reg %v", msgType))
	}
	msgCbs[msgType] = cb
}
