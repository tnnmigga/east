package pm

import (
	"github.com/tnnmigga/nett/idef"
	"github.com/tnnmigga/nett/msgbus"
)

var manager *Manager

type Manager struct {
	idef.IModule
	cache   map[uint64]*Player
	waiting map[uint64][]func(*Player, error)
}

func Init(m idef.IModule) {
	manager = &Manager{
		IModule: m,
		cache:   map[uint64]*Player{},
		waiting: map[uint64][]func(*Player, error){},
	}
	msgbus.RegisterHandler(m, manager.onC2SPackage)
	msgbus.RegisterRPC(m, manager.onCreatePlayerRPC)
}
