package eventimpl

import (
	"east/core/eventbus"
	"east/game/compts/play/domain"
	"east/game/compts/play/domain/api"
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
