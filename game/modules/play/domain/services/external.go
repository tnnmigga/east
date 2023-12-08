package services

import (
	"east/game/modules/play/domain"
	"east/game/modules/play/domain/services/eventimpl"
	"east/game/modules/play/domain/services/msgimpl"
	"east/game/modules/play/domain/services/userimpl"
)

func Init(d *domain.Domain) {
	d.PutImpl(domain.MsgImplIndex, msgimpl.New(d))
	d.PutImpl(domain.UserImplIndex, userimpl.New(d))
	d.PutImpl(domain.EventImplIndex, eventimpl.New(d))
	d.Init()
}
