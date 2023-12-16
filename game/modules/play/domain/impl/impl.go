package impl

import (
	"east/game/modules/play/domain"
	"east/game/modules/play/domain/impl/eventimpl"
	"east/game/modules/play/domain/impl/msgimpl"
	"east/game/modules/play/domain/impl/timerimpl"
	"east/game/modules/play/domain/impl/userimpl"
)

func Init(d *domain.Domain) {
	d.PutImpl(domain.MsgImplIndex, msgimpl.New(d))
	d.PutImpl(domain.EventImplIndex, eventimpl.New(d))
	d.PutImpl(domain.TimerImplIndex, timerimpl.New(d))
	d.PutImpl(domain.UserImplIndex, userimpl.New(d))
	d.Init()
}
