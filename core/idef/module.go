package idef

import "reflect"

type IModule interface {
	Name() string
	MQ() chan any
	Handlers() map[reflect.Type]*Handler
	Run()
	Stop()
}
