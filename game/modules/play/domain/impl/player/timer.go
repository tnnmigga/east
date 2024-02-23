package player

import (
	"east/core/codec"
	"east/core/msgbus"
	"east/pb"
)

type timerSayHello struct {
	Text string
}

func (c *useCase) onTimerSayHello(ctx *timerSayHello) {
	msgbus.Cast(&pb.S2CPackage{
		UserID: 1,
		Body: codec.Encode(&pb.SayHelloResp{
			Text: ctx.Text,
		}),
	}, msgbus.NonuseStream(), msgbus.ServerID(1888))
}
