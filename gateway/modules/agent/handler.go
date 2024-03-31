package agent

import (
	"east/core/log"
	"east/core/msgbus"
	"east/pb"
)

func (m *module) initHandler() {
	msgbus.RegisterHandler(m, m.onS2CPackage)
	msgbus.RegisterRPC(m, m.onTestRPC)
}

func (m *module) onS2CPackage(pkg *pb.S2CPackage) {
	agent := m.manager.GetAgent(pkg.UserID)
	if agent == nil {
		log.Warnf("agent not found %d", pkg.UserID)
		return
	}
	select {
	case agent.sendMQ <- pkg.Body:
	default:
		log.Errorf("agent send mq full! %d", pkg.UserID)
	}
}

func (m *module) onTestRPC(pkg *pb.TestRPC, resolve func(any), reject func(error)) {
	resolve(&pb.TestRPCRes{
		V: 11,
	})
}
