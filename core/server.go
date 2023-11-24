package core

import (
	"eden/core/log"
	"eden/core/module"
	"eden/core/pb"
	"eden/core/util"
	"sync"
)

type Server struct {
	modules []module.IModule
	mq      chan *pb.Package
}

func (s *Server) Run(modules ...module.IModule) (stopFn func()) {
	wg := &sync.WaitGroup{}
	for _, m := range modules {
		wg.Add(1)
		go m.Run(wg)
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
