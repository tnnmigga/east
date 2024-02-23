package player

import (
	"east/core/idef"
	"east/core/msgbus"
	"east/define"
	"east/game/modules/play/domain"
	"east/game/modules/play/domain/api"
)

type service struct {
	*domain.Domain
}

func New(d *domain.Domain) api.IMsg {
	s := &service{
		Domain: d,
	}
	s.After(idef.ServerStateInit, s.afterInit)
	return s
}

func (s *service) afterInit() error {
	msgbus.RegisterHandler(s, s.onSayHelloReq)
	msgbus.RegisterHandler(s, s.onTimerSayHello)
	msgbus.RegisterRPC(s, s.onRPCTest)
	s.EventCase().RegisterHandler(define.EventUserSayHello, s.onEventUserMsg)
	return nil
}
