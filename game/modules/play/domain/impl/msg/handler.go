package msg

import (
	"east/core/codec"
	"east/core/log"
	"east/core/msgbus"
	"east/core/util"
	"east/pb"
)

func (c *useCase) onC2SPackage(msg *pb.C2SPackage) {
	pkg, err := codec.Decode(msg.Body)
	if err != nil {
		log.Errorf("msg decode error %v", err)
		return
	}
	log.Errorf("recv user %d msg %s", msg.UserID, util.String(pkg))
	msgbus.CastLocal(pkg)
}
