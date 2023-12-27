package core

import (
	"east/core/iconf"
	"east/core/idef"
	"east/core/infra"
	"east/core/infra/nats"
	"east/core/log"
	"east/core/sys"
	"east/core/util"
	"fmt"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

type Server struct {
	modules []idef.IModule
	wg      *sync.WaitGroup
}

func NewServer(modules ...idef.IModule) idef.IServer {
	server := &Server{
		modules: make([]idef.IModule, 0, len(modules)+1),
		wg:      &sync.WaitGroup{},
	}
	server.modules = append(server.modules, nats.New(infra.ModTypNats)) // nats最后停止
	server.modules = append(server.modules, modules...)
	return server
}

func (s *Server) Init() {
	s.before(idef.ServerStateInit, s.abort)
	iconf.LoadFromJSON(util.ReadFile("configs.jsonc"))
	log.Init()
	log.Info("server initialization")
	s.after(idef.ServerStateInit, s.abort)
}

func (s *Server) Run() {
	s.before(idef.ServerStateRun, s.abort)
	log.Info("server try to run")
	for _, m := range s.modules {
		s.runModule(s.wg, m)
	}
	log.Info("server running successfully")
	s.after(idef.ServerStateRun, s.abort)
}

func (s *Server) Stop() {
	s.before(idef.ServerStateStop, s.noabort)
	log.Info("server try to stop")
	sys.WaitGoDone(time.Minute)
	for i := len(s.modules) - 1; i >= 0; i-- {
		m := s.modules[i]
		util.ExecAndRecover(m.Stop)
	}
	s.wg.Wait()
	log.Info("server stoped successfully")
	s.after(idef.ServerStateStop, s.noabort)
}

func (s *Server) Close() {
	s.before(idef.ServerStateClose, s.noabort)
	log.Info("server close")
	os.Exit(0)
}

func (s *Server) runModule(wg *sync.WaitGroup, m idef.IModule) {
	wg.Add(1)
	go func() {
		defer util.RecoverPanic()
		defer wg.Done()
		m.Run()
	}()
}

func (s *Server) abort(m idef.IModule, err error) {
	log.Fatalf("module %s, on %s, error: %v", m.Name(), util.Caller(3), err)
}

func (s *Server) noabort(m idef.IModule, err error) {
	log.Errorf("module %s, on %s, error: %v", m.Name(), util.Caller(3), err)
}

func (s *Server) before(state idef.ServerState, onError ...func(idef.IModule, error)) {
	for _, m := range s.modules {
		hook := m.Hook(state, 0)
		for _, h := range hook {
			if err := wrapHook(h)(); err != nil {
				log.Errorf("server before %#v error, module %s, error %v", state, m.Name(), err)
				for _, f := range onError {
					f(m, err)
				}
			}
		}
	}
}

func (s *Server) after(state idef.ServerState, onError ...func(idef.IModule, error)) {
	for _, m := range s.modules {
		hook := m.Hook(state, 1)
		for _, h := range hook {
			if err := wrapHook(h)(); err != nil {
				log.Errorf("server before %#v error, module %s, error %v", state, m.Name(), err)
				for _, f := range onError {
					f(m, err)
				}
			}
		}
	}
}

// wrapHook 添加panic处理
func wrapHook(h func() error) func() error {
	return func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("%v: %s", r, debug.Stack())
			}
		}()
		return h()
	}
}
