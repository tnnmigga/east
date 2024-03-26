package agent

import (
	"context"
	"east/core/basic"
	"east/core/conf"
	"east/core/idef"
	"east/core/log"
	"east/define"
	"net"
	"sync"
)

type module struct {
	*basic.Module
	sync.RWMutex
	lister net.Listener
	conns  map[uint64]*userAgent
}

func New() idef.IModule {
	m := &module{
		Module: basic.New(define.ModTypTCPAgent, basic.DefaultMQLen),
		conns:  map[uint64]*userAgent{},
	}
	m.initHandler()
	m.After(idef.ServerStateInit, m.afterInit)
	m.After(idef.ServerStateRun, m.afterRun)
	return m
}

func (m *module) afterInit() (err error) {
	m.lister, err = net.Listen("tcp", conf.String("gateway.addr", "127.0.0.1:9527"))
	if err != nil {
		return err
	}
	return nil
}

func (m *module) afterRun() error {
	go m.accept()
	return nil
}

func (m *module) accept() {
	for {
		conn, err := m.lister.Accept()
		log.Infof("new conn!")
		if err != nil {
			log.Errorf("tcpagent accept error %v", err)
			continue
		}
		uid := uint64(1) //idgen.NewUUID()
		ctx, cancel := context.WithCancel(context.Background())
		agent := &userAgent{
			ctx:      ctx,
			userInfo: userInfo{userID: uid, serverID: 1},
			conn:     conn,
			mq:       make(chan []byte, 256),
			Destory:  cancel,
		}
		m.Lock()
		m.conns[uid] = agent
		m.Unlock()
		agent.run()
	}
}
