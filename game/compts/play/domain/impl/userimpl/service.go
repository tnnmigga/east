package userimpl

import (
	define1 "east/core/idef"
	"east/core/msgbus"
	"east/define"
	"east/game/compts/play/domain"
	"east/game/compts/play/domain/api"
)

type service struct {
	*domain.Domain
}

func New(d *domain.Domain) api.IMsg {
	s := &service{
		Domain: d,
	}
	s.After(define1.ServerStateInit, s.afterInit)
	return s
}

func (s *service) afterInit() error {
	msgbus.RegisterHandler(s, s.onSayHelloReq)
	msgbus.RegisterHandler(s, s.onTimerSayHello)
	s.EventImpl().RegisterHandler(define.EventUserSayHello, s.onEventUserMsg)
	return nil
}
