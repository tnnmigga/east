package userimpl

import (
	"east/core/codec"
	"east/core/msgbus"
	"east/define"
	"east/game/modules/play/domain/services/userimpl/usermeta"
	"east/pb"
)

func (s *service) onTimerSayHello(ctx *usermeta.TimerSayHello) {
	msgbus.Broadcast(define.ServTypGateway, &pb.S2CPackage{
		UserID: 1,
		Body: codec.Encode(&pb.SayHelloResp{
			Text: ctx.Text,
		}),
	})
}
