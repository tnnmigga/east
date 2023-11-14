package core

import (
	"eden/common/algorithms"
	"fmt"
	"reflect"

	"github.com/gogo/protobuf/proto"
)

type MsgCb func(msg proto.Message)

var msgCbs = map[reflect.Type]MsgCb{}

var msgIDtoType = map[uint32]reflect.Type{}


func RegMsgCb(msg proto.Message, cb MsgCb) {
	msgType := reflect.TypeOf(msg)
	_, ok := msgCbs[msgType]
	if ok {
		panic(fmt.Errorf("RegMsgCb duplicated reg %v", msgType))
	}
	msgCbs[msgType] = cb
	msgid := MsgID(msg)
	if _, ok := msgIDtoType[msgid]; ok {
		panic(fmt.Errorf("RegMsgCb msgid clash %d", msgid))
	}
	msgIDtoType[msgid] = msgType
}

func MsgID(v any) uint32 {
	name := reflect.TypeOf(v).Elem().Name()
	return algorithms.BKDRHash([]byte(name))
}
