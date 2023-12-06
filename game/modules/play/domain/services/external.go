package services

import (
	"east/game/modules/play/domain"
	"east/game/modules/play/domain/services/msgimpl"
	"east/game/modules/play/domain/services/userimpl"
)

func Init(d *domain.Domain) {
	d.PutCase(domain.MsgCaseIndex, msgimpl.New(d))
	d.PutCase(domain.UserCaseIndex, userimpl.New(d))
	d.Init()
}
