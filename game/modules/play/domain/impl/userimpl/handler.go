package userimpl

import (
	"east/core/eventbus"
	"east/core/iconf"
	"east/core/log"
	"east/core/msgbus"
	"east/define"
	"east/pb"
	"time"
)

func (s *service) registerMsgHandler() {
	msgbus.RegisterHandler(s, s.onSayHelloReq)
}

func (s *service) onSayHelloReq(msg *pb.SayHelloReq) {
	log.Infof("client say hello %v", msg.Text)
	msgbus.Cast(iconf.ServerID(), &eventbus.Event{
		OwnerID: 1,
		Topic:   define.EventUserSayHello,
		Value:   1,
	})
	s.TimerImpl().Create(time.Second*2, &timerSayHello{
		Text: "hello client!",
	})
}
