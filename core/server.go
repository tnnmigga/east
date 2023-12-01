package core

import (
	"east/core/define"
	"east/core/infra"
	"east/core/infra/nats"
	"east/core/log"
	"east/core/util"
	"sync"
)

type Server struct {
	modules []define.IModule
	// mq      chan any // 对内分发
}

func NewServer(modules ...define.IModule) *Server {
	return &Server{
		modules: append(modules, nats.New(infra.ModTypNats)), // nats module最后启动最先停止
	}
}

func (s *Server) Run() (stop func()) {
	wg := &sync.WaitGroup{}
	for _, m := range s.modules {
		s.NewGoroutine(wg, m.Run)
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

func (s *Server) NewGoroutine(wg *sync.WaitGroup, fn func()) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer util.RecoverPanic()
		fn()
	}()
}
