package idef

import (
	"reflect"
)

type IModule interface {
	Name() string
	MQ() chan any
	Run()
	Stop()
	RegisterHandler(mType reflect.Type, handler *Handler)
	Hook(state ServerState, stage int) []func() error
}
