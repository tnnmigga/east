package play

import (
	"east/define"
	"east/game/modules/play/domain"
	"east/game/modules/play/domain/impl"

	"github.com/tnnmigga/nett/basic"
	"github.com/tnnmigga/nett/idef"
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
