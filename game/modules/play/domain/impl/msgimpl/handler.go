package msgimpl

import (
	"east/core/codec"
	"east/core/iconf"
	"east/core/log"
	"east/core/msgbus"
	"east/core/util"
	"east/pb"
)

func (m *service) onC2SPackage(msg *pb.C2SPackage) {
	pkg, err := codec.Decode(msg.Body)
	if err != nil {
		log.Errorf("msg decode error %v", err)
		return
	}
	log.Errorf("recv user %d msg %s", msg.UserID, util.String(pkg))
	msgbus.Cast(iconf.ServerID(), pkg)
}
