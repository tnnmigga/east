package play

import (
	"east/define"
	"east/game/play/domain"
	"east/game/play/domain/impl"
	"east/game/play/pm"

	"github.com/tnnmigga/core/idef"
	"github.com/tnnmigga/core/mods/basic"
)

type module struct {
	*basic.Module
	*domain.Domain
}

func New() idef.IModule {
	m := &module{
		Module: basic.New(define.ModPlay, basic.DefaultMQLen),
	}
	m.Domain = domain.New(m)
	pm.Init(m)
	impl.Init(m.Domain)
	return m
}
