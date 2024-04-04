package player

import (
	"east/core/eventbus"
	"east/core/msgbus"
	"east/core/zlog"
	"east/define"
	"east/pb"
	"time"
)

func (c *useCase) onSayHelloReq(msg *pb.SayHelloReq) {
	zlog.Infof("client say hello %v", msg.Text)
	msgbus.CastLocal(&eventbus.Event{
		OwnerID: 1,
		Topic:   define.EventUserSayHello,
		Value:   1,
	})
	c.TimerCase().New(time.Second*2, &timerSayHello{
		UserID: 1,
		Text:   "hello client!",
	})
}

func (c *useCase) onRPCTest(req *pb.TestRPC, resolve func(any), reject func(error)) {
	resolve(&pb.TestRPCRes{
		V: 22,
	})
}
