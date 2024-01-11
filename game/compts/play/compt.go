package play

import (
	"east/core/compt"
	"east/core/idef"
	"east/define"
	"east/game/compts/play/domain"
	"east/game/compts/play/domain/impl"
)

type component struct {
	*domain.Domain
}

func New() idef.IComponent {
	com := &component{
		Domain: domain.New(compt.New(define.ModTypPlay, compt.DefaultMQLen)),
	}
	impl.Init(com.Domain)
	return com
}
