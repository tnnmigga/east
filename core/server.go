package core

import (
	"east/core/idef"
	"east/core/infra"
	"east/core/infra/nats"
	"east/core/log"
	"east/core/sys"
	"east/core/utils"
	"sync"
	"time"
)

type Server struct {
	modules []idef.IModule
	// mq      chan any // 对内分发
}

func NewServer(modules ...idef.IModule) *Server {
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
		log.Infof("try stop modules...")
		sys.WaitGoDone(time.Minute)
		for i := len(s.modules) - 1; i >= 0; i-- {
			m := s.modules[i]
			utils.ExecAndRecover(m.Stop)
		}
		wg.Wait()
		log.Infof("stop modules success")
	}
}

func (s *Server) NewGoroutine(wg *sync.WaitGroup, fn func()) {
	wg.Add(1)
	go func() {
		defer utils.RecoverPanic()
		defer wg.Done()
		fn()
	}()
}
