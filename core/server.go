package core

import (
	"eden/core/iconf"
	"eden/core/infra"
	"eden/core/infra/nats"
	"eden/core/log"
	"eden/core/message"
	"eden/core/module"
	"eden/core/util"
	"sync"
)

type Server struct {
	modules []module.IModule
	mq      chan *message.Package // 对内分发
}

func NewServer() *Server {
	return &Server{
		mq: make(chan *message.Package, iconf.Int32("mq-len")),
	}
}

func (s *Server) Run(modules ...module.IModule) (stopFn func()) {
	wg := &sync.WaitGroup{}
	s.modules = make([]module.IModule, 0, len(modules)+1)
	s.modules = append(s.modules, nats.NewModule(infra.Nats))
	s.modules = append(s.modules, modules...)
	message.Attach(s.mq)
	s.modules = append(s.modules, modules...)
	for _, m := range s.modules {
		wg.Add(1)
		go s.run(wg, m.Run)
	}
	return func() {
		log.Infof("try close server, start close modules")
		for i := len(s.modules) - 1; i >= 0; i-- {
			m := s.modules[i]
			util.ExecAndRecover(m.Close)
		}
		wg.Wait()
		log.Infof("close modules success")
	}
}

func (s *Server) dispatch() {
	for pkg := range s.mq {
		for _, m := range s.modules {
			if pkg.Module != m.Name() {
				continue
			}
			select {
			case m.MQ() <- pkg:
			default:
				log.Errorf("server dispatch mq full %v", m.Name())
				if pkg.TTL < iconf.Int32("msg-max-ttl", 1) { // 消息堆积
					pkg.TTL++
					s.mq <- pkg
				}
			}
		}
	}
}

func (s *Server) run(wg *sync.WaitGroup, fn func()) {
	defer wg.Done()
	defer util.RecoverPanic()
	fn()
}
