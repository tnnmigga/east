package impl

import (
	"east/game/compts/play/domain"
	"east/game/compts/play/domain/impl/eventimpl"
	"east/game/compts/play/domain/impl/msgimpl"
	"east/game/compts/play/domain/impl/timerimpl"
	"east/game/compts/play/domain/impl/userimpl"
)

func Init(d *domain.Domain) {
	d.PutImpl(domain.MsgImplIndex, msgimpl.New(d))
	d.PutImpl(domain.EventImplIndex, eventimpl.New(d))
	d.PutImpl(domain.TimerImplIndex, timerimpl.New(d))
	d.PutImpl(domain.UserImplIndex, userimpl.New(d))
}
