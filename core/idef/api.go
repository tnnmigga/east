package idef

import (
	"reflect"
)

type IComponent interface {
	Name() string
	Assign(any)
	MQ() chan any
	Run()
	Stop()
	RegisterHandler(mType reflect.Type, handler *Handler)
	Hook(state ServerState, stage int) []func() error
}
