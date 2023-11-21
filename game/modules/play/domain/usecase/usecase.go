package usecase

import (
	"eden/game/modules/play/domain"
	"eden/game/modules/play/domain/usecase/msgcase"
)

func Init(d *domain.Domain) {
	d.PutCase(domain.MsgCaseIndex, msgcase.NewCase(d))
}
