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
	m.RLock()
	defer m.RUnlock()
	agent, ok := m.conns[pkg.UserID]
	if !ok {
		return
	}
	select {
	case agent.mq <- pkg.Body:
	default:
		log.Errorf("userAgent mq full!")
		agent.close()
	}
}

func (m *module) onTestRPC(pkg *pb.TestRPC, resolve func(any), reject func(error)) {
	resolve(&pb.TestRPCRes{
		V: 11,
	})
}
