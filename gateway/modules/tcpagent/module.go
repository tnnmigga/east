package tcpagent

import (
	"east/core/iconf"
	"east/core/idef"
	"east/core/module"
	"east/core/util/idgen"
	define1 "east/define"
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
	lister, err := net.Listen("tcp", iconf.String("tcp-addr", "127.0.0.1:9527"))
	if err != nil {
		panic(err)
	}
	m := &Module{
		Module: module.New(define1.ModTypTCPAgent, 100000),
		lister: lister,
		conns:  map[uint64]*userAgent{},
	}
	m.initHandler()
	return m
}

func (m *Module) Run() {
	go m.accept()
	m.Module.Run()
}

func (m *Module) accept() {
	for {
		conn, err := m.lister.Accept()
		if err != nil {
			continue
		}
		uid := idgen.NewUUID()
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
