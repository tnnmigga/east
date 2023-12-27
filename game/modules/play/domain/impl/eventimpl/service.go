package eventimpl

import (
	"east/core/eventbus"
	"east/core/idef"
	"east/game/modules/play/domain"
	"east/game/modules/play/domain/api"
)

type service struct {
	*domain.Domain
	*eventbus.EventBus
}

func New(d *domain.Domain) api.IEvent {
	s := &service{
		Domain:   d,
		EventBus: eventbus.New(d),
	}
	s.After(idef.ServerStateInit, s.afterInit)
	s.Before(idef.ServerStateStop, s.beforeStop)
	return s
}

func (s *service) afterInit() error {
	// 加载数据
	return nil
}

func (s *service) beforeStop() error {
	// 数据落地
	return nil
}
