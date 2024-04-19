package agent

import (
	"east/pb"

	"github.com/tnnmigga/nett/infra/zlog"
	"github.com/tnnmigga/nett/msgbus"
)

func (m *module) initHandler() {
	msgbus.RegisterHandler(m, m.onS2CPackage)
	msgbus.RegisterRPC(m, m.onTestRPC)
}

func (m *module) onS2CPackage(pkg *pb.S2CPackage) {
	agent := m.manager.GetAgent(pkg.UserID)
	if agent == nil {
		zlog.Warnf("agent not found %d", pkg.UserID)
		return
	}
	select {
	case agent.sendMQ <- pkg.Body:
	default:
		zlog.Errorf("agent send mq full! %d", pkg.UserID)
	}
}

func (m *module) onTestRPC(pkg *pb.TestRPC, resolve func(any), reject func(error)) {
	resolve(&pb.TestRPCRes{
		V: 11,
	})
}
