package userimpl

import (
	"east/game/modules/play/domain"
	"east/game/modules/play/domain/api"
)

type service struct {
	*domain.Domain
}

func New(d *domain.Domain) api.IMsg {
	return &service{}
}

func (s *service) Init() {
}

func (s *service) Destroy() {
}
