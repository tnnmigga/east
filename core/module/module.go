package module

import (
	"eden/core/log"
	"eden/core/util"
	"reflect"

	"github.com/gogo/protobuf/proto"
)

type ModuleType string

type IModule interface {
	Type() ModuleType
	MQ() chan<- proto.Message
	Init()
	Run()
	Close()
}

type Module struct {
	typ      ModuleType
	mq       chan proto.Message
	handlers map[reflect.Type]MessageHandler
}

var modules = map[int64]*Module{}

func NewModule(mType ModuleType, mqLen int32) *Module {
	m := &Module{
		typ:      mType,
		mq:       make(chan proto.Message, mqLen),
		handlers: map[reflect.Type]MessageHandler{},
	}
	return m
}

func (m *Module) Type() ModuleType {
	return m.typ
}

func (m *Module) MQ() chan<- proto.Message {
	return m.mq
}

func (m *Module) Init() {
}

func (m *Module) Run() {
	for msg := range m.mq {
		msgType := reflect.TypeOf(msg)
		cb, ok := m.handlers[msgType]
		if !ok {
			log.Errorf("handler not exist %v", msgType)
			continue
		}
		func() {
			defer util.PrintPanicStack()
			cb(msg)
		}()
	}
}

func (m *Module) Close() {
	close(m.mq)
}
