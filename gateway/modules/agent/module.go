package agent

import (
	"east/define"
	"fmt"

	"github.com/tnnmigga/nett/basic"
	"github.com/tnnmigga/nett/idef"
)

type AgentType string

const (
	AgentTypeTCP       AgentType = "tcp"
	AgentTypeWebSocket AgentType = "websocket"
)

type module struct {
	*basic.Module
	agentType AgentType
	listener  IListener
	manager   *AgentManager
}

type IListener interface {
	Run()
	Close()
}

func New(agentType AgentType) idef.IModule {
	m := &module{
		Module:    basic.New(define.ModAgent, basic.DefaultMQLen),
		agentType: agentType,
	}
	m.manager = NewAgentManager(m)
	m.initHandler()
	m.After(idef.ServerStateInit, m.afterInit)
	m.After(idef.ServerStateRun, m.afterRun)
	m.Before(idef.ServerStateStop, m.beforeStop)
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
	m.listener.Run()
	return nil
}

func (m *module) beforeStop() error {
	m.listener.Close()
	return nil
}
