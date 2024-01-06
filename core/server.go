package core

import (
	"east/core/conf"
	"east/core/idef"
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
	compts []idef.IComponent
	wg     *sync.WaitGroup
}

func NewServer(compts ...idef.IComponent) idef.IServer {
	server := &Server{
		compts: make([]idef.IComponent, 0, len(compts)+1),
		wg:     &sync.WaitGroup{},
	}
	server.compts = append(server.compts, nats.New()) // nats最后停止
	server.compts = append(server.compts, compts...)
	return server
}

func (s *Server) Init() {
	s.before(idef.ServerStateInit, s.abort)
	conf.LoadFromJSON(util.ReadFile("configs.jsonc"))
	log.Init()
	log.Info("server initialization")
	s.after(idef.ServerStateInit, s.abort)
}

func (s *Server) Run() {
	s.before(idef.ServerStateRun, s.abort)
	log.Info("server try to run")
	for _, com := range s.compts {
		s.runCompt(s.wg, com)
	}
	log.Info("server running successfully")
	s.after(idef.ServerStateRun, s.abort)
}

func (s *Server) Stop() {
	s.before(idef.ServerStateStop, s.noabort)
	log.Info("server try to stop")
	s.waitMsgHandling(time.Minute)
	sys.WaitGoDone(time.Minute)
	for i := len(s.compts) - 1; i >= 0; i-- {
		com := s.compts[i]
		util.ExecAndRecover(com.Stop)
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

func (s *Server) runCompt(wg *sync.WaitGroup, com idef.IComponent) {
	wg.Add(1)
	go func() {
		defer util.RecoverPanic()
		defer wg.Done()
		com.Run()
	}()
}

func (s *Server) waitMsgHandling(maxWaitTime time.Duration) {
	maxCheckCount := maxWaitTime / time.Second * 10
	for ; maxCheckCount > 0; maxCheckCount-- {
		time.Sleep(100 * time.Millisecond)
		isEmpty := true
		for _, com := range s.compts {
			if len(com.MQ()) != 0 {
				isEmpty = false
				break
			}
		}
		if isEmpty {
			return
		}
	}
	log.Errorf("wait msg handing timeout")
}

func (s *Server) abort(com idef.IComponent, err error) {
	log.Fatalf("component %s, on %s, error: %v", com.Name(), util.Caller(3), err)
}

func (s *Server) noabort(com idef.IComponent, err error) {
	log.Errorf("component %s, on %s, error: %v", com.Name(), util.Caller(3), err)
}

func (s *Server) before(state idef.ServerState, onError ...func(idef.IComponent, error)) {
	for _, com := range s.compts {
		hook := com.Hook(state, 0)
		for _, h := range hook {
			if err := wrapHook(h)(); err != nil {
				log.Errorf("server before %#v error, component %s, error %v", state, com.Name(), err)
				for _, f := range onError {
					f(com, err)
				}
			}
		}
	}
}

func (s *Server) after(state idef.ServerState, onError ...func(idef.IComponent, error)) {
	for _, com := range s.compts {
		hook := com.Hook(state, 1)
		for _, h := range hook {
			if err := wrapHook(h)(); err != nil {
				log.Errorf("server before %#v error, component %s, error %v", state, com.Name(), err)
				for _, f := range onError {
					f(com, err)
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
