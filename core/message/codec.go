package message

import (
	"eden/core/algorithms"
	"reflect"

	"github.com/gogo/protobuf/proto"
)

func Encode(v proto.Message) []byte {
	return nil
}

func MsgID(v any) uint32 {
	name := reflect.TypeOf(v).Elem().Name()
	return algorithms.BKDRHash([]byte(name))
}
