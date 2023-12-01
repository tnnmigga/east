package usecase

import (
	"east/game/modules/play/domain"
	"east/game/modules/play/domain/usecase/msgcase"
)

func Init(d *domain.Domain) {
	d.PutCase(domain.MsgCaseIndex, msgcase.NewCase(d))
}
