package idef

import (
	"reflect"
)

type IModule interface {
	Name() string
	Assign(any)
	Run()
	Stop()
	RegisterHandler(mType reflect.Type, handler *Handler)
	Hook(state ServerState, stage int) []func() error
}

type IServer interface {
	Init()
	Run()
	Stop()
	Close()
}
