package tcpagent

import (
	"east/core/iconf"
	"east/core/idef"
	"east/core/log"
	"east/core/module"
	"east/define"
	"net"
	"sync"
)

type Module struct {
	*module.Module
	sync.RWMutex
	lister net.Listener
	conns  map[uint64]*userAgent
}

func New() idef.IModule {
	m := &Module{
		Module: module.New(define.ModTypTCPAgent, module.DefaultMQLen),
		conns:  map[uint64]*userAgent{},
	}
	m.initHandler()
	m.After(idef.ServerStateInit, m.afterInit)
	m.After(idef.ServerStateRun, m.afterRun)
	return m
}

func (m *Module) afterInit() (err error) {
	m.lister, err = net.Listen("tcp", iconf.String("tcp-addr", "127.0.0.1:9527"))
	if err != nil {
		return err
	}
	return nil
}

func (m *Module) afterRun() error {
	go m.accept()
	return nil
}

func (m *Module) accept() {
	for {
		conn, err := m.lister.Accept()
		log.Infof("new conn!")
		if err != nil {
			log.Errorf("tcpagent accept error %v", err)
			continue
		}
		uid := uint64(1) //idgen.NewUUID()
		agent := &userAgent{
			userInfo: userInfo{userID: uid, serverID: 1},
			conn:     conn,
			mq:       make(chan []byte, 256),
		}
		m.Lock()
		m.conns[uid] = agent
		m.Unlock()
		agent.run()
	}
}
