package msgimpl

import (
	"east/core/log"
	"east/core/message"
	"east/define"
	"east/pb"
)

func (s *service) regMsgHandler() {
	message.RegisterHandler(s, s.onSayHelloReq)
}

func (m *service) onSayHelloReq(msg *pb.SayHelloReq) {
	log.Infof("onSayHelloReq %v", msg.Text)
	message.Broadcast(define.ServTypGateway, &pb.SayHelloResp{
		Text: "hello, client",
	})
}
