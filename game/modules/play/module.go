package play

import (
	"east/core/idef"
	"east/core/module"
	"east/define"
	"east/game/modules/play/domain"
	"east/game/modules/play/domain/services"
)

type Module struct {
	*domain.Domain
}

func New() idef.IModule {
	basicModule := module.New(define.ModTypPlay, 100000)
	m := &Module{
		Domain: domain.NewDomain(basicModule),
	}
	services.Init(m.Domain)
	return m
}
