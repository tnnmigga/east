package services

import (
	"east/game/modules/play/domain"
	"east/game/modules/play/domain/services/msgimpl"
)

func Init(d *domain.Domain) {
	d.PutCase(domain.MsgCaseIndex, msgimpl.New(d))
}
