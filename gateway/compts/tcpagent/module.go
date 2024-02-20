package tcpagent

import (
	"east/core/compt"
	"east/core/conf"
	define1 "east/core/idef"
	"east/core/log"
	"east/define"
	"net"
	"sync"
)

type component struct {
	*compt.Module
	sync.RWMutex
	lister net.Listener
	conns  map[uint64]*userAgent
}

func New() define1.IComponent {
	com := &component{
		Module: compt.New(define.ModTypTCPAgent, compt.DefaultMQLen),
		conns:  map[uint64]*userAgent{},
	}
	com.initHandler()
	com.After(define1.ServerStateInit, com.afterInit)
	com.After(define1.ServerStateRun, com.afterRun)
	return com
}

func (com *component) afterInit() (err error) {
	com.lister, err = net.Listen("tcp", conf.String("tcp-addr", "127.0.0.1:9527"))
	if err != nil {
		return err
	}
	return nil
}

func (com *component) afterRun() error {
	go com.accept()
	return nil
}

func (com *component) accept() {
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
