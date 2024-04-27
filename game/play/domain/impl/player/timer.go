package player

import (
	"east/pb"

	"github.com/tnnmigga/nett/codec"
	"github.com/tnnmigga/nett/msgbus"
)

type timerSayHello struct {
	UserID uint64
	Text   string
}

func (c *useCase) onTimerSayHello(ctx *timerSayHello) {
	msgbus.Cast(&pb.S2CPackage{
		UserID: ctx.UserID,
		Body: codec.Encode(&pb.SayHelloResp{
			Text: ctx.Text,
		}),
	}, msgbus.UseStream(), msgbus.ServerID(1888))
}
