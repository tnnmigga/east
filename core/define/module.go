package define

import "reflect"

type IModule interface {
	Name() string
	MQ() chan any
	Handlers() map[reflect.Type]*HandlerFn
	Run()
	Close()
}