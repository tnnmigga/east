package core

import (
	"eden/core/health"
	"eden/core/log"
	"eden/core/module"
	"eden/core/util"
	"eden/pb"
	"sync"
)

type Server struct {
	modules []module.IModule
	mq      chan *pb.Package
}

func (s *Server) Run(modules ...module.IModule) {
	wg := &sync.WaitGroup{}
	for _, m := range modules {
		wg.Add(1)
		m.Init()
		go m.Run(wg)
	}
	sign := health.WaitExitSignal()
	log.Infof("receive exit signal %v, try close server", sign)
	s.stop()
	wg.Wait()
	log.Info("close server success")
}

func (s *Server) stop() {
	for i := len(s.modules) - 1; i >= 0; i-- {
		m := s.modules[i]
		defer func() {
			util.PrintPanicStack()
			m.Close()
		}()
	}
}
