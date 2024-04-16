package codec

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/tnnmigga/nett/util"
	"github.com/tnnmigga/nett/zlog"

	"github.com/gogo/protobuf/proto"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	msgIDToDesc map[uint32]*MessageDescriptor
)

func init() {
	msgIDToDesc = map[uint32]*MessageDescriptor{}
}

const (
	marshalTypeGogoproto = iota
	marshalTypeBSON
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
	name := util.TypeName(v)
	id := util.TypeID(v)
	if desc, has := msgIDToDesc[id]; has {
		if desc.MessageName != name {
			zlog.Fatalf("msgid duplicat %v %d", name, id)
		}
	}
	mType := reflect.TypeOf(v)
	if mType.Kind() == reflect.Ptr {
		mType = mType.Elem()
	}
	msgIDToDesc[id] = &MessageDescriptor{
		MessageName: name,
		MarshalType: marshalType(v),
		ReflectType: mType,
	}
}

func Encode(v any) []byte {
	msgID := util.TypeID(v)
	bytes := Marshal(v)
	body := make([]byte, 4, len(bytes)+4)
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
		return nil, fmt.Errorf("message decode msgid not found %d", msgID)
	}
	msg = desc.New()
	err = Unmarshal(b[4:], msg)
	return msg, err
}

func Marshal(v any) []byte {
	if v0, ok := v.(proto.Message); ok {
		b, err := proto.Marshal(v0)
		if err != nil {
			log.Panic(fmt.Errorf("message encode error %v", err))
		}
		return b
	}
	b, err := bson.Marshal(v)
	if err != nil {
		log.Panic(fmt.Errorf("message encode error %v", err))
	}
	return b
}

func Unmarshal(b []byte, addr any) error {
	switch marshalType(addr) {
	case marshalTypeGogoproto:
		return proto.Unmarshal(b, addr.(proto.Message))
	case marshalTypeBSON:
		return bson.Unmarshal(b, addr)
	default:
		return errors.New("invalid marshal type")
	}
}

func marshalType(v any) int {
	if _, ok := v.(proto.Message); ok {
		return marshalTypeGogoproto
	}
	return marshalTypeBSON
}
