package agent

import (
	"context"
	"east/core/basic"
	"east/core/log"
	"east/core/msgbus"
	"east/core/util"
	"east/pb"
	"sync"
	"sync/atomic"
	"time"
)

const (
	MaxSendMQLen = 64
)

const (
	AgentStateNone = iota
	AgentStateRun
	AgentStateWait
	AgentStateClose
)

type AgentManager struct {
	agents map[uint64]*Agent
	rw     sync.RWMutex
}

func NewAgentManager() *AgentManager {
	return &AgentManager{
		agents: map[uint64]*Agent{},
	}
}

func (am *AgentManager) AddAgent(agent *Agent) {
	if agent.userID == 0 {
		log.Errorf("agent uid is 0")
		return
	}
	am.rw.Lock()
	am.agents[agent.userID] = agent
	am.rw.Unlock()
}

func (am *AgentManager) RemoveAgent(uid uint64) {
	am.rw.Lock()
	delete(am.agents, uid)
	am.rw.Unlock()
}

func (am *AgentManager) GetAgent(uid uint64) *Agent {
	am.rw.RLock()
	agent := am.agents[uid]
	am.rw.RUnlock()
	return agent
}

func (am *AgentManager) OnConnect(conn Conn) {
	agent := NewAgent()
	agent.conn = conn
	conn.BindAgent(agent)
	agent.Run()
}

func (am *AgentManager) OnError(err error) {
	log.Errorf("agent manager error %v", err)
}

type IAgent interface {
	OnMessage([]byte)
	OnClose()
	OnReconnect()
	OnReadError(error)
}

type Agent struct {
	userID  uint64
	servID  uint32
	ctx     context.Context
	conn    Conn
	sendMQ  chan []byte
	closeFn func()
	state   int32
	failMsg []byte
	waitMs  time.Duration
}

func NewAgent() *Agent {
	agent := &Agent{
		sendMQ: make(chan []byte, MaxSendMQLen),
	}
	agent.ctx, agent.closeFn = context.WithCancel(context.Background())
	return agent
}

func (a *Agent) OnMessage(data []byte) {
	msgbus.Cast(&pb.C2SPackage{
		UserID: a.userID,
		Body:   data,
	}, msgbus.ServerID(a.servID))
}

func (a *Agent) OnClose() {
	a.closeFn()
}

func (a *Agent) OnReadError(err error) {
	if atomic.CompareAndSwapInt32(&a.state, AgentStateRun, AgentStateWait) {
		a.waitMs = util.NowNs()
		log.Debugf("agent read error %v", err)
		a.conn.Close()
	}
}

func (a *Agent) OnWriteError(data []byte, err error) {
	if atomic.CompareAndSwapInt32(&a.state, AgentStateRun, AgentStateWait) {
		a.waitMs = util.NowNs()
		a.conn.Close()
	}
	a.failMsg = data
}

func (a *Agent) OnReconnect() {
	if atomic.CompareAndSwapInt32(&a.state, AgentStateWait, AgentStateRun) {
		log.Debugf("agent reconnect")
	}
	if a.failMsg != nil {
		if err := a.conn.Write(a.failMsg); err != nil {
			a.OnWriteError(a.failMsg, err)
			return
		}
	}
	a.Run()
}

func (a *Agent) Run() {
	basic.Go(a.writeLoop)
	basic.Go(func() {
		a.conn.Run(a.ctx)
	})
}

func (a *Agent) writeLoop() {
	for data := range a.sendMQ {
		if atomic.LoadInt32(&a.state) != AgentStateRun {
			a.failMsg = data
			return
		}
		if util.ContextDone(a.ctx) {
			return
		}
		err := a.conn.Write(data)
		if err != nil {
			a.OnWriteError(data, err)
			return
		}
	}
}

type Conn interface {
	Run(context.Context)
	Write([]byte) error
	Close()
	BindAgent(IAgent)
}
