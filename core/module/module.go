package module

import (
	"eden/core/log"
	"eden/core/util"
	"reflect"
	"sync"

	"github.com/gogo/protobuf/proto"
)

type ModuleName string

type IModule interface {
	Name() ModuleName
	MQ() chan<- proto.Message
	Init()
	Run(wg *sync.WaitGroup)
	Close()
}

type Module struct {
	name     ModuleName
	mq       chan proto.Message
	handlers map[reflect.Type]MessageHandler
}

var modules = map[int64]*Module{}

func NewModule(mType ModuleName, mqLen int32) *Module {
	m := &Module{
		name:     mType,
		mq:       make(chan proto.Message, mqLen),
		handlers: map[reflect.Type]MessageHandler{},
	}
	return m
}

func (m *Module) Name() ModuleName {
	return m.name
}

func (m *Module) MQ() chan<- proto.Message {
	return m.mq
}

func (m *Module) Init() {
}

func (m *Module) Run(wg *sync.WaitGroup) {
	defer func() {
		log.Infof("%v module exist success", m.name)
		wg.Done()
	}()
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
