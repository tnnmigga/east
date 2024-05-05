package event

import (
	"east/game/play/domain"
	"east/game/play/domain/api"

	"github.com/tnnmigga/core/infra/eventbus"
)

type useCase struct {
	*domain.Domain
	*eventbus.EventBus
}

func New(d *domain.Domain) api.IEvent {
	return &useCase{
		Domain:   d,
		EventBus: eventbus.New(d),
	}
}
