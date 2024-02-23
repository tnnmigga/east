package player

import (
	"east/core/idef"
	"east/core/msgbus"
	"east/define"
	"east/game/modules/play/domain"
	"east/game/modules/play/domain/api"
)

type useCase struct {
	*domain.Domain
}

func New(d *domain.Domain) api.IMsg {
	s := &useCase{
		Domain: d,
	}
	s.After(idef.ServerStateInit, s.afterInit)
	return s
}

func (c *useCase) afterInit() error {
	msgbus.RegisterHandler(c, c.onSayHelloReq)
	msgbus.RegisterHandler(c, c.onTimerSayHello)
	msgbus.RegisterRPC(c, c.onRPCTest)
	c.EventCase().RegisterHandler(define.EventUserSayHello, c.onEventUserMsg)
	return nil
}
