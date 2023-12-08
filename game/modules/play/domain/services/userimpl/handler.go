package userimpl

import (
	"east/core/codec"
	"east/core/eventbus"
	"east/core/iconf"
	"east/core/log"
	"east/core/message"
	"east/define"
	"east/pb"
)

func (s *service) onSayHelloReq(msg *pb.SayHelloReq) {
	log.Infof("onSayHelloReq %v", msg.Text)
	message.Broadcast(define.ServTypGateway, &pb.S2CPackage{
		UserID: 1,
		Body: codec.Encode(&pb.SayHelloResp{
			Text: "hello, client!",
		}),
	})
	message.Cast(iconf.ServerID(), &eventbus.Event{
		OwnerID: 1,
		Topic:   define.EventUserSayHello,
		Value:   1,
	})
}
