package play

import (
	"east/core/idef"
	"east/core/module"
	"east/define"
	"east/game/modules/play/domain"
	"east/game/modules/play/domain/services"
)

type Module struct {
	*module.Module
	*domain.Domain
}

func NewModule() idef.IModule {
	m := &Module{
		Module: module.New(define.ModTypPlay, 100000),
		Domain: domain.NewDomain(),
	}
	services.Init(m.Domain)
	return m
}
