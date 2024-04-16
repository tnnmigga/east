package nett

import (
	"fmt"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"github.com/tnnmigga/nett/conf"
	"github.com/tnnmigga/nett/core"
	"github.com/tnnmigga/nett/idef"
	"github.com/tnnmigga/nett/infra/link"
	"github.com/tnnmigga/nett/util"
	"github.com/tnnmigga/nett/zlog"
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
	server.modules = append(server.modules, link.New()) // nats最后停止
	server.modules = append(server.modules, modules...)
	server.init()
	server.run()
	return server
}

func (s *Server) init() {
	conf.LoadFromJSON(util.ReadFile("configs.jsonc"))
	zlog.Init()
	zlog.Info("server initialization")
	s.after(idef.ServerStateInit, s.abort)
}

func (s *Server) run() {
	s.before(idef.ServerStateRun, s.abort)
	zlog.Info("server try to run")
	for _, m := range s.modules {
		s.runModule(s.wg, m)
	}
	zlog.Info("server running successfully")
	s.after(idef.ServerStateRun, s.abort)
}

func (s *Server) stop() {
	s.before(idef.ServerStateStop, s.noabort)
	zlog.Info("server try to stop")
	s.waitMsgHandling(time.Minute)
	core.WaitGoDone(time.Minute)
	for i := len(s.modules) - 1; i >= 0; i-- {
		m := s.modules[i]
		util.ExecAndRecover(m.Stop)
	}
	s.wg.Wait()
	zlog.Info("server stoped successfully")
	s.exit()
}

func (s *Server) Exit() {
	s.stop()
}

func (s *Server) exit() {
	s.before(idef.ServerStateExit, s.noabort)
	zlog.Info("server close")
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

func (s *Server) waitMsgHandling(maxWaitTime time.Duration) {
	// 每100ms检查一次模块消息是否处理完
	maxCheckCount := maxWaitTime / time.Millisecond / 100
	for ; maxCheckCount > 0; maxCheckCount-- {
		time.Sleep(100 * time.Millisecond)
		isEmpty := true
		for _, m := range s.modules {
			if len(m.MQ()) != 0 {
				isEmpty = false
				break
			}
		}
		if isEmpty {
			return
		}
	}
	zlog.Errorf("wait msg handing timeout")
}

func (s *Server) abort(m idef.IModule, err error) {
	zlog.Fatalf("module %s, on %s, error: %v", m.Name(), util.Caller(3), err)
}

func (s *Server) noabort(m idef.IModule, err error) {
	zlog.Errorf("module %s, on %s, error: %v", m.Name(), util.Caller(3), err)
}

func (s *Server) before(state idef.ServerState, onError ...func(idef.IModule, error)) {
	for _, m := range s.modules {
		hook := m.Hook(state, 0)
		for _, h := range hook {
			if err := wrapHook(h)(); err != nil {
				zlog.Errorf("server before %#v error, module %s, error %v", state, m.Name(), err)
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
				zlog.Errorf("server before %#v error, module %s, error %v", state, m.Name(), err)
				for _, f := range onError {
					f(m, err)
				}
			}
		}
	}
}

// 添加panic处理
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
