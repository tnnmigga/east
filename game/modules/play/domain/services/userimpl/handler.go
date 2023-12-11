package userimpl

import (
	"east/core/eventbus"
	"east/core/iconf"
	"east/core/log"
	"east/core/msgbus"
	"east/define"
	"east/game/modules/play/domain/services/userimpl/usermeta"
	"east/pb"
	"time"
)

func (s *service) onSayHelloReq(msg *pb.SayHelloReq) {
	log.Infof("onSayHelloReq %v", msg.Text)
	msgbus.Cast(iconf.ServerID(), &eventbus.Event{
		OwnerID: 1,
		Topic:   define.EventUserSayHello,
		Value:   1,
	})
	s.TimerImpl().After(time.Second*2, &usermeta.TimerSayHello{
		Text: "hello client!",
	})
}
