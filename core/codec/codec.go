package codec

import (
	"eden/core/util"
	"eden/core/util/idgen"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"

	"github.com/gogo/protobuf/proto"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	msgIDToDesc map[uint32]*MessageDescriptor
)

const (
	marshalTypeGogoproto = 1
	marshalTypeBSON      = 2
)

type MessageDescriptor struct {
	MessageName string
	MarshalType int
	ReflectType reflect.Type
}

func (d *MessageDescriptor) New() any {
	return reflect.New(d.ReflectType).Interface()
}

func Register(v any) {
	name := util.ReflectName(v)
	id := msgID(v)
	if _, has := msgIDToDesc[id]; has {
		panic(fmt.Errorf("msgid duplicat %v %d", name, id))
	}
	if _, ok := v.(proto.Message); ok {
		msgIDToDesc[id] = &MessageDescriptor{
			MessageName: name,
			MarshalType: marshalType(v),
			ReflectType: reflect.TypeOf(v),
		}
	}
}

func Encode(v any) []byte {
	msgID := msgID(v)
	bytes := toBytes(v)
	body := make([]byte, len(bytes)+4)
	binary.LittleEndian.PutUint32(body, msgID)
	body = append(body, bytes...)
	return body
}

func Decode(b []byte) (msg any, err error) {
	if len(b) < 4 {
		return nil, fmt.Errorf("message decode len error %d", len(b))
	}
	msgID := binary.LittleEndian.Uint32(b)
	desc, ok := msgIDToDesc[msgID]
	if !ok {
		return nil, fmt.Errorf("message decode msgid error %d", msgID)
	}
	msg = desc.New()
	switch desc.MarshalType {
	case marshalTypeGogoproto:
		err = proto.Unmarshal(b, msg.(proto.Message))
	case marshalTypeBSON:
		err = bson.Unmarshal(b, msg)
	default:
		err = errors.New("invalid marshal type")
	}
	return msg, err
}

func toBytes(v any) []byte {
	if v0, ok := v.(proto.Message); ok {
		b, err := proto.Marshal(v0)
		if err != nil {
			panic(fmt.Errorf("message encode error %v", err))
		}
		return b
	}
	b, err := bson.Marshal(v)
	if err != nil {
		panic(fmt.Errorf("message encode error %v", err))
	}
	return b
}

func msgID(v any) uint32 {
	name := util.ReflectName(v)
	return idgen.HashToID(name)
}

func marshalType(v any) int {
	if _, ok := v.(proto.Message); ok {
		return marshalTypeGogoproto
	}
	return marshalTypeBSON
}