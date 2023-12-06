package msgimpl

import (
	"east/core/codec"
	"east/core/iconf"
	"east/core/log"
	"east/core/message"
	"east/pb"
)

func (s *service) regMsgHandler() {
	message.RegisterHandler(s, s.onC2SPackage)
}

func (m *service) onC2SPackage(msg *pb.C2SPackage) {
	log.Infof("recv client msg %v", msg.String())
	pkg, err := codec.Decode(msg.Body)
	if err != nil {
		log.Errorf("onC2SPackage decode error %v", err)
		return
	}
	message.Cast(iconf.ServerID(), pkg)
}
