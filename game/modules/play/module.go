package play

import (
	"east/core/define"
	"east/core/module"
	define1 "east/define"
	"east/game/modules/play/domain"
	"east/game/modules/play/domain/usecase"
)

type Module struct {
	*module.Module
	*domain.Domain
}

func NewModule() define.IModule {
	m := &Module{
		Module: module.New(define1.ModTypPlay, 100000),
		Domain: domain.NewDomain(),
	}
	usecase.Init(m.Domain)
	return m
}
