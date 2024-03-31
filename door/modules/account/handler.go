package account

import (
	"east/core/msgbus"
	"east/pb"
)

func (m *module) initHandler() {
	msgbus.RegisterRPC(m, m.onTokenAuthReq)
}

func (m *module) onTokenAuthReq(req *pb.TokenAuthReq, resolve func(any), reject func(error)) {
	resolve(&pb.TokenAuthResp{
		Code:    pb.SUCC,
		UserID:  1, // util.RandomInterval[uint64](1, 1000),
		SeverID: 1999,
	})
}
