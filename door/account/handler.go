package account

import (
	"east/define"
	"east/pb"

	"github.com/tnnmigga/core/infra/zlog"
	"github.com/tnnmigga/core/msgbus"
)

func (m *module) initHandler() {
	msgbus.RegisterRPC(m, m.onTokenAuthReq)
}

func (m *module) onTokenAuthReq(req *pb.TokenAuthReq, resolve func(any), reject func(error)) {
	zlog.Debugf("onTokenAuthReq: %v", req)
	msgbus.RPC(m, msgbus.ServerType(define.ServGame), &pb.CreatePlayerRPC{
		UserID: 1,
	}, func(res *pb.CreatePlayerRPCRes, err error) {
		if err != nil {
			reject(err)
			return
		}
		resolve(&pb.TokenAuthResp{
			Code:    pb.SUCCESS,
			UserID:  1,
			SeverID: res.ServerID,
		})
	})
}
