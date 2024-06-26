package account

import (
	"east/define"

	"github.com/tnnmigga/core/conf"
	"github.com/tnnmigga/core/idef"
	"github.com/tnnmigga/core/infra/https"
	"github.com/tnnmigga/core/mods/basic"
)

type module struct {
	*basic.Module
	agent *https.HttpAgent
}

func New() idef.IModule {
	m := &module{
		Module: basic.New(define.ModAccount, basic.DefaultMQLen),
	}
	m.After(idef.ServerStateInit, m.afterInit)
	m.After(idef.ServerStateRun, m.afterRun)
	m.Before(idef.ServerStateStop, m.beforeStop)
	return m
}

func (m *module) afterInit() error {
	m.initHandler()
	m.initRoute()
	return nil
}

func (m *module) afterRun() error {
	err := m.agent.Run(conf.String("account.addr", "127.0.0.1:8080"))
	return err
}

func (m *module) beforeStop() error {
	err := m.agent.Stop()
	return err
}
