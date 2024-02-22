package event

import (
	"east/core/eventbus"
	"east/game/modules/play/domain"
	"east/game/modules/play/domain/api"
)

type service struct {
	*domain.Domain
	*eventbus.EventBus
}

func New(d *domain.Domain) api.IEvent {
	return &service{
		Domain:   d,
		EventBus: eventbus.New(d),
	}
}
