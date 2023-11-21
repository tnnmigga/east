package msgcase

import (
	"eden/game/modules/play/domain"
	"eden/game/modules/play/domain/api"
)

type useCase struct {
}

func NewCase(d *domain.Domain) api.IMsgCase {
	return &useCase{}
}

func (uc *useCase) Init() {
}

func (uc *useCase) Destroy() {
}
