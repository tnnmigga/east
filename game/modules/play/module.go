package play

import (
	"east/core/idef"
	"east/core/module"
	"east/define"
	"east/game/modules/play/domain"
	"east/game/modules/play/domain/impl"
)

type Module struct {
	*domain.Domain
}

func New() idef.IModule {
	m := &Module{
		Domain: domain.New(module.New(define.ModTypPlay, module.DefaultMQLen)),
	}
	impl.Init(m.Domain)
	return m
}
