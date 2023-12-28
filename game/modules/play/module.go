package play

import (
	"east/core/idef"
	"east/core/mod"
	"east/define"
	"east/game/modules/play/domain"
	"east/game/modules/play/domain/impl"
)

type module struct {
	*domain.Domain
}

func New() idef.IModule {
	m := &module{
		Domain: domain.New(mod.New(define.ModTypPlay, mod.DefaultMQLen)),
	}
	impl.Init(m.Domain)
	return m
}
