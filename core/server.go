package core

import (
	"eden/core/module"
	"fmt"
)

type Server struct {
	modules map[module.ModuleType]module.IModule
}

func (s *Server) Start(modules ...module.IModule) {
	for _, m := range modules {
		mType := m.Type()
		if _, exist := s.modules[mType]; exist {
			panic(fmt.Errorf("module type repead %v", mType))
		}
		
	}
}
