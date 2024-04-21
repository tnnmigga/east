package player

import (
	"east/define"
	"east/game/play/domain"
	"east/game/play/domain/api"
	"east/game/play/pm"

	"github.com/tnnmigga/nett/idef"
	"github.com/tnnmigga/nett/msgbus"
)

type useCase struct {
	*domain.Domain
}

func New(d *domain.Domain) api.IMsg {
	c := &useCase{
		Domain: d,
	}
	c.After(idef.ServerStateInit, c.afterInit)
	return c
}

func (c *useCase) afterInit() error {
	msgbus.RegisterHandler(c, c.onTimerSayHello)
	pm.RegMsgHandler(c, c.onSayHelloReq)
	msgbus.RegisterRPC(c, c.onRPCTest)
	c.EventCase().RegisterHandler(define.EventUserSayHello, c.onEventUserMsg)
	return nil
}
