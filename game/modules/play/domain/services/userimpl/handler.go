package userimpl

import (
	"east/core/log"
	"east/core/message"
	"east/pb"
)

func (s *service) regMsgHandler() {
	message.RegisterHandler(s, s.onSayHelloReq)
}

func (m *service) onSayHelloReq(msg *pb.SayHelloReq) {
	log.Infof("onSayHelloReq %v", msg.Text)
}
