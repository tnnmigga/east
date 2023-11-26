package module

import (
	"eden/core/log"
	"eden/core/util"
	"reflect"
)

type IModule interface {
	Name() string
	MQ() chan<- any
	Run()
	Close()
	Handlers() map[reflect.Type]MessageHandler
}

type Module struct {
	name     string
	mq       chan any
	handlers map[reflect.Type]MessageHandler
}

var modules = map[int64]*Module{}

func NewModule(mType string, mqLen int32) *Module {
	m := &Module{
		name:     mType,
		mq:       make(chan any, mqLen),
		handlers: map[reflect.Type]MessageHandler{},
	}
	return m
}

func (m *Module) Name() string {
	return m.name
}

func (m *Module) MQ() chan<- any {
	return m.mq
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
			defer util.RecoverPanic()
			cb(msg)
		}()
	}
	log.Infof("%v module exist success", m.name)
}

func (m *Module) Close() {
	close(m.mq)
}

func (m *Module) Handlers() map[reflect.Type]MessageHandler {
	return m.handlers
}
