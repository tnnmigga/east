package msg

import (
	"east/pb"

	"github.com/tnnmigga/nett/codec"
	"github.com/tnnmigga/nett/msgbus"
	"github.com/tnnmigga/nett/util"
	"github.com/tnnmigga/nett/zlog"
)

func (c *useCase) onC2SPackage(msg *pb.C2SPackage) {
	pkg, err := codec.Decode(msg.Body)
	if err != nil {
		zlog.Errorf("msg decode error %v", err)
		return
	}
	zlog.Errorf("recv user %d msg %s", msg.UserID, util.String(pkg))
	msgbus.Cast(pkg)
}
