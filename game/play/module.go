package play

import (
	"east/define"
	"east/game/play/domain"
	"east/game/play/domain/impl"

	"github.com/tnnmigga/nett/idef"
	"github.com/tnnmigga/nett/mods/basic"
)

type module struct {
	*domain.Domain
}

func New() idef.IModule {
	m := &module{
		Domain: domain.New(basic.New(define.ModPlay, basic.DefaultMQLen)),
	}
	impl.Init(m.Domain)
	return m
}
