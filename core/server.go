package core

import (
	"eden/core/configs"
	"eden/core/infra/nats"
	"eden/core/log"
	"eden/core/message"
	"eden/core/module"
	"eden/core/util"
	"runtime/debug"
	"sync"
)

type Server struct {
	modules []module.IModule
	mq      chan *message.Package // 对内分发
}

func NewServer() *Server {
	return &Server{
		mq: make(chan *message.Package, configs.Int32("mq-len")),
	}
}

func (s *Server) Run(modules ...module.IModule) (stopFn func()) {
	wg := &sync.WaitGroup{}
	s.modules = append(s.modules, nats.NewModule())
	message.Attach(s.mq, s.modules[0].MQ())
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
		if pkg.ServerID != configs.ServerID() {

		}
	}
}

func (s *Server) run(wg *sync.WaitGroup, fn func()) {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("%v: %s", r, debug.Stack())
		}
		wg.Done()
	}()
	fn()
}
