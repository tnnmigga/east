package nats

import (
	"eden/core/module"

	"github.com/nats-io/nats.go"
)

type Module struct {
	*module.Module
	conn *nats.Conn
}

func NewModule(name module.ModuleName, url string) module.IModule {
	conn, err := nats.Connect(url)
	if err != nil {
		panic(err)
	}
	m := &Module{
		Module: module.NewModule(name, 100000),
		conn:   conn,
	}
	return m
}

func (m *Module) recv() {

}
