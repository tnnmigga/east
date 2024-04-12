package agent

import (
	"context"
	"east/define"
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
	MaxAliveTime = 10 * time.Second
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
	am := &AgentManager{
		IModule: m,
		agents:  map[uint64]*Agent{},
	}
	core.Go(am.clean)
	return am
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
	msgbus.RPC(am, msgbus.ServerType(define.ServDoor), &pb.TokenAuthReq{}, func(resp *pb.TokenAuthResp, err error) {
		if err != nil {
			zlog.Errorf("TokenAuthReq return error %v", err)
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

func (am *AgentManager) clean(ctx context.Context) {
	ticker := time.NewTicker(MaxAliveTime)
	for {
		select {
		case <-ticker.C:
			am.cleanDeadAgent()
		case <-ctx.Done():
			return
		}
	}
}

func (am *AgentManager) cleanDeadAgent() {
	am.rw.RLock()
	defer am.rw.RUnlock()
	nowNs := util.NowNs()
	for uid, agent := range am.agents {
		state := atomic.LoadInt32(&agent.state)
		if state == AgentStateDead {
			delete(am.agents, uid)
			continue
		}
		if state != AgentStateWait {
			continue
		}
		if agent.waitNs+MaxAliveTime > nowNs {
			continue
		}
		if !atomic.CompareAndSwapInt32(&agent.state, AgentStateWait, AgentStateDead) {
			continue
		}
		delete(am.agents, uid)
	}
}

type IAgent interface {
	OnMessage([]byte)
	OnReconnect()
	OnReadError(error)
}

type Agent struct {
	userID  uint64
	servID  uint32
	conn    Conn
	sendMQ  chan []byte
	state   int32
	failMsg []byte
	waitNs  time.Duration
	cancel  func()
}

func NewAgent() *Agent {
	agent := &Agent{
		sendMQ: make(chan []byte, MaxSendMQLen),
	}
	return agent
}

func (a *Agent) OnMessage(data []byte) {
	msgbus.Cast(&pb.C2SPackage{
		UserID: a.userID,
		Body:   data,
	}, msgbus.ServerID(a.servID))
}

func (a *Agent) OnReadError(err error) {
	if errors.Is(err, io.EOF) {
		atomic.StoreInt32(&a.state, AgentStateDead)
		a.conn.Close()
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
	ctx, cancel := context.WithCancel(context.Background())
	if a.cancel != nil {
		zlog.Warnf("cancel is exist")
	}
	a.cancel = cancel
	a.conn.ReadLoop(ctx)
	a.writeLoop(ctx)
}

func (a *Agent) writeLoop(ctx context.Context) {
	core.Go(func() {
		for {
			select {
			case data := <-a.sendMQ:
				if atomic.LoadInt32(&a.state) != AgentStateRun {
					a.failMsg = data
					return
				}
				err := a.conn.Write(data)
				if err == nil {
					continue
				}
				a.OnWriteError(data, err)
				return
			case <-ctx.Done():
				return
			}
		}
	})
}

func (a *Agent) beforeStop() {
	a.cancel()
}

func (a *Agent) afterStop() {
	a.conn.Close()
}

type Conn interface {
	ReadLoop(ctx context.Context)
	Write([]byte) error
	Close()
	BindAgent(IAgent)
}
