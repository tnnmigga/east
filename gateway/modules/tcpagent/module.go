package tcpagent

import (
	"east/core/basic"
	"east/core/conf"
	define1 "east/core/idef"
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

func New() define1.IModule {
	com := &module{
		Module: basic.New(define.ModTypTCPAgent, basic.DefaultMQLen),
		conns:  map[uint64]*userAgent{},
	}
	com.initHandler()
	com.After(define1.ServerStateInit, com.afterInit)
	com.After(define1.ServerStateRun, com.afterRun)
	return com
}

func (com *module) afterInit() (err error) {
	com.lister, err = net.Listen("tcp", conf.String("tcp-addr", "127.0.0.1:9527"))
	if err != nil {
		return err
	}
	return nil
}

func (com *module) afterRun() error {
	go com.accept()
	return nil
}

func (com *module) accept() {
	for {
		conn, err := com.lister.Accept()
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
		com.Lock()
		com.conns[uid] = agent
		com.Unlock()
		agent.run()
	}
}
