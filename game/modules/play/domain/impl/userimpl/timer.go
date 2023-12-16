package userimpl

import (
	"east/core/codec"
	"east/core/msgbus"
	"east/pb"
)

type timerSayHello struct {
	Text string
}

func (s *service) onTimerSayHello(ctx *timerSayHello) {
	// msgbus.Cast(1888, &pb.S2CPackage{
	// 	UserID: 1,
	// 	Body: codec.Encode(&pb.SayHelloResp{
	// 		Text: ctx.Text,
	// 	}),
	// })
	msgbus.Cast(1888, &pb.S2CPackage{
		UserID: 1,
		Body: codec.Encode(&pb.SayHelloResp{
			Text: ctx.Text,
		}),
	}, msgbus.NonuseStream())
}
