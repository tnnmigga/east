package play

import (
	"east/core/com"
	define1 "east/core/idef"
	"east/define"
	"east/game/modules/play/domain"
	"east/game/modules/play/domain/impl"
)

type module struct {
	*domain.Domain
}

func New() define1.IModule {
	m := &module{
		Domain: domain.New(com.New(define.ModTypPlay, com.DefaultMQLen)),
	}
	impl.Init(m.Domain)
	return m
}
