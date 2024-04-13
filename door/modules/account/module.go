package account

import (
	"east/define"

	"github.com/tnnmigga/nett/conf"
	"github.com/tnnmigga/nett/idef"
	"github.com/tnnmigga/nett/modules/basic"
	"github.com/tnnmigga/nett/web"
)

type module struct {
	*basic.Module
	agent *web.HttpAgent
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
