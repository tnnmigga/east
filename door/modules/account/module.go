package account

import (
	"east/core/basic"
	"east/core/conf"
	"east/core/idef"
	"east/core/web"
	"east/define"
)

type module struct {
	*basic.Module
	agent *web.HttpAgent
}

func New() idef.IModule {
	m := &module{
		Module: basic.New(define.ModTypAccount, basic.DefaultMQLen),
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
