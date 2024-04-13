package msg

import (
	"east/pb"

	"github.com/tnnmigga/nett/codec"
	"github.com/tnnmigga/nett/infra/zlog"
	"github.com/tnnmigga/nett/msgbus"
	"github.com/tnnmigga/nett/util"
)

func (c *useCase) onC2SPackage(req *pb.C2SPackage) {
	pkg, err := codec.Decode(req.Body)
	if err != nil {
		zlog.Errorf("msg decode error %v", err)
		return
	}
	zlog.Errorf("recv user %d msg %s", req.UserID, util.String(pkg))
	msgbus.Cast(pkg)
}
