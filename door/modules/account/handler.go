package account

import (
	"east/pb"

	"github.com/tnnmigga/nett/msgbus"
	"github.com/tnnmigga/nett/zlog"
)

func (m *module) initHandler() {
	msgbus.RegisterRPC(m, m.onTokenAuthReq)
}

func (m *module) onTokenAuthReq(req *pb.TokenAuthReq, resolve func(any), reject func(error)) {
	zlog.Debugf("onTokenAuthReq: %v", req)
	resolve(&pb.TokenAuthResp{
		Code:    pb.SUCCESS,
		UserID:  1, // util.RandomInterval[uint64](1, 1000),
		SeverID: 1999,
	})
}
