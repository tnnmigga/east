package play

import (
	"east/core/basic"
	"east/core/idef"
	"east/define"
	"east/game/modules/play/domain"
	"east/game/modules/play/domain/impl"
)

type module struct {
	*domain.Domain
}

func New() idef.IModule {
	m := &module{
		Domain: domain.New(basic.New(define.ModTypPlay, basic.DefaultMQLen)),
	}
	impl.Init(m.Domain)
	return m
}
