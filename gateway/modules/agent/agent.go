package agent

import (
	"context"
	"east/pb"
	"errors"
	"io"
	"runtime"
	"sync"

	"github.com/tnnmigga/nett/core"
	"github.com/tnnmigga/nett/idef"
	"github.com/tnnmigga/nett/msgbus"
	"github.com/tnnmigga/nett/util"
	"github.com/tnnmigga/nett/zlog"

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
	AgentStateDead
)

type AgentManager struct {
	idef.IModule
	agents map[uint64]*Agent
	rw     sync.RWMutex
}

func NewAgentManager(m idef.IModule) *AgentManager {
	return &AgentManager{
		IModule: m,
		agents:  map[uint64]*Agent{},
	}
}

func (am *AgentManager) AddAgent(agent *Agent) {
	if agent.userID == 0 {
		zlog.Errorf("agent uid is 0")
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
	msgbus.RPC(am, 4242, &pb.TokenAuthReq{}, func(resp *pb.TokenAuthResp, err error) {
		if err != nil {
			conn.Close()
			return
		}
		agent := NewAgent()
		agent.conn = conn
		agent.servID = resp.SeverID
		agent.userID = resp.UserID
		conn.BindAgent(agent)
		am.AddAgent(agent)
		agent.Run()
	})

}

func (am *AgentManager) OnError(err error) {
	zlog.Errorf("agent manager error %v", err)
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
	waitNs  time.Duration
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
	if errors.Is(err, io.EOF) {
		atomic.StoreInt32(&a.state, AgentStateDead)
	}
	if atomic.CompareAndSwapInt32(&a.state, AgentStateRun, AgentStateWait) {
		a.waitNs = util.NowNs()
		zlog.Debugf("agent read error %v", err)
		a.conn.Close()
	}
	runtime.Goexit()
}

func (a *Agent) OnWriteError(data []byte, err error) {
	if atomic.CompareAndSwapInt32(&a.state, AgentStateRun, AgentStateWait) {
		a.waitNs = util.NowNs()
		a.conn.Close()
	}
	a.failMsg = data
	runtime.Goexit()
}

func (a *Agent) OnReconnect() {
	if atomic.CompareAndSwapInt32(&a.state, AgentStateWait, AgentStateRun) {
		zlog.Debugf("agent reconnect")
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
	a.state = AgentStateRun
	// 启动agent读写goroutine
	core.Go(a.writeLoop)
	core.Go(func() {
		a.conn.Run(a.ctx)
	})
}

func (a *Agent) writeLoop() {
	for data := range a.sendMQ {
		if atomic.LoadInt32(&a.state) != AgentStateRun {
			a.failMsg = data
			return
		}
		err := a.conn.Write(data)
		if err != nil {
			a.OnWriteError(data, err)
		}
	}
}

type Conn interface {
	Run(context.Context)
	Write([]byte) error
	Close()
	BindAgent(IAgent)
}
