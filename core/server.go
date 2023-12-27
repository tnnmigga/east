package core

import (
	"east/core/iconf"
	"east/core/idef"
	"east/core/infra"
	"east/core/infra/nats"
	"east/core/log"
	"east/core/sys"
	"east/core/util"
	"os"
	"sync"
	"time"
)

type Server struct {
	modules []idef.IModule
	wg      *sync.WaitGroup
}

func NewServer(modules ...idef.IModule) *Server {
	server := &Server{
		modules: make([]idef.IModule, 0, len(modules)+1),
		wg:      &sync.WaitGroup{},
	}
	server.modules = append(server.modules, nats.New(infra.ModTypNats)) // nats最后停止
	server.modules = append(server.modules, modules...)
	return server
}

func (s *Server) Init() {
	s.before(idef.ServerStateInit, s.panicErr)
	iconf.LoadFromJSON(util.ReadFile("configs.jsonc"))
	log.Init()
	s.after(idef.ServerStateInit)
}

func (s *Server) Run() {
	s.before(idef.ServerStateRun, s.panicErr)
	for _, m := range s.modules {
		s.runModule(s.wg, m)
	}
	s.after(idef.ServerStateRun, s.panicErr)
}

func (s *Server) Stop() {
	s.before(idef.ServerStateStop, s.logErr)
	sys.WaitGoDone(time.Minute)
	for i := len(s.modules) - 1; i >= 0; i-- {
		m := s.modules[i]
		util.ExecAndRecover(m.Stop)
	}
	s.wg.Wait()
	s.after(idef.ServerStateStop, s.logErr)
}

func (s *Server) Close() {
	s.before(idef.ServerStateClose, s.logErr)
	os.Exit(0)
}

func (s *Server) panicErr(err error) {
	panic(err)
}

func (s *Server) logErr(err error) {
	log.Fatal(err)
}

func (s *Server) before(state idef.ServerState, onError ...func(error)) {
	for _, m := range s.modules {
		hook := m.Hook(state, 0)
		for _, h := range hook {
			if err := h(); err != nil {
				log.Errorf("server before %#v error, module %s, error %v", state, m.Name(), err)
				for _, f := range onError {
					f(err)
				}
			}
		}
	}
}

func (s *Server) after(state idef.ServerState, onError ...func(error)) {
	for _, m := range s.modules {
		hook := m.Hook(state, 1)
		for _, h := range hook {
			if err := h(); err != nil {
				log.Errorf("server before %#v error, module %s, error %v", state, m.Name(), err)
				for _, f := range onError {
					f(err)
				}
			}
		}
	}
}

func (s *Server) runModule(wg *sync.WaitGroup, m idef.IModule) {
	wg.Add(1)
	go func() {
		defer util.RecoverPanic()
		defer wg.Done()
		m.Run()
	}()
}
