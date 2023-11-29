package play

import (
	"eden/core/module"
	"eden/game/modules/play/domain"
	"eden/game/modules/play/domain/usecase"
)

type Module struct {
	*module.Module
	*domain.Domain
}

func NewModule() module.IModule {
	m := &Module{
		Module: module.New("play", 100000),
		Domain: domain.NewDomain(),
	}
	usecase.Init(m.Domain)
	return m
}
