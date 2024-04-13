package event

import (
	"east/game/modules/play/domain"
	"east/game/modules/play/domain/api"

	"github.com/tnnmigga/nett/infra/eventbus"
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
