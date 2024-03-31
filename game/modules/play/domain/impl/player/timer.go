package player

import (
	"east/core/codec"
	"east/core/msgbus"
	"east/pb"
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
	}, msgbus.NonuseStream(), msgbus.ServerID(1888))
}
