package core

import (
	"eden/core/module"
	"eden/core/signal"
	"eden/pb"
)

type Server struct {
	modules []module.IModule
	mq      chan *pb.Package
}

func (s *Server) Run(modules ...module.IModule) {
	for _, m := range modules {
		m.Init()
		go m.Run()
	}
	signal.WaitStopSignal()
	for i := len(s.modules) - 1; i >= 0; i-- {
		m := s.modules[i]
		m.Close()
	}
}
