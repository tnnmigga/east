package agent

import (
	"east/define"
	"fmt"

	"github.com/tnnmigga/nett/idef"
	"github.com/tnnmigga/nett/mods/basic"
)

const (
	AgentTypeTCP       = "tcp"
	AgentTypeWebSocket = "websocket"
)

type module struct {
	*basic.Module
	agentType string
	listener  IListener
	manager   *AgentManager
}

type IListener interface {
	Run() error
	Close()
}

func New(agentType string) idef.IModule {
	m := &module{
		Module:    basic.New(define.ModAgent, basic.DefaultMQLen),
		agentType: agentType,
	}
	m.manager = NewAgentManager(m)
	m.initHandler()
	m.After(idef.ServerStateInit, m.afterInit)
	m.After(idef.ServerStateRun, m.afterRun)
	m.Before(idef.ServerStateStop, m.beforeStop)
	m.After(idef.ServerStateStop, m.afterStop)
	return m
}

func (m *module) afterInit() (err error) {
	switch m.agentType {
	case AgentTypeTCP:
		m.listener = NewTCPListener(m.manager)
	case AgentTypeWebSocket:
		m.listener = NewWebSocketListener(m.manager)
	default:
		return fmt.Errorf("unknown agent type: %s", m.agentType)
	}
	return nil
}

func (m *module) afterRun() error {
	return m.listener.Run()
}

func (m *module) beforeStop() error {
	m.listener.Close()
	for _, agent := range m.manager.agents {
		agent.beforeStop()
	}
	return nil
}

func (m *module) afterStop() error {
	for _, agent := range m.manager.agents {
		agent.afterStop()
	}
	return nil
}
