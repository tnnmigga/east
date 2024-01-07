package userimpl

import (
	"east/core/eventbus"
	"east/core/log"
	"east/core/msgbus"
	"east/define"
	"east/pb"
	"time"
)

func (s *service) onSayHelloReq(msg *pb.SayHelloReq) {
	log.Infof("client say hello %v", msg.Text)
	msgbus.CastLocal(&eventbus.Event{
		OwnerID: 1,
		Topic:   define.EventUserSayHello,
		Value:   1,
	})
	s.TimerImpl().New(time.Second*2, &timerSayHello{
		Text: "hello client!",
	})
}

func (s *service) onRPCTest(req *pb.TestRPC, resolve func(any), reject func(error)) {
	resolve(&pb.TestRPCRes{
		V: 22,
	})
}
