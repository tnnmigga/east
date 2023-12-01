package tcpagent

import (
	"east/core/log"
	"east/core/message"
	"east/pb"
)

func (m *Module) initHandler() {
	message.RegisterHandler(m, m.onS2CPackage)
}

func (m *Module) onS2CPackage(pkg *pb.S2CPackage) {
	m.RLock()
	agent, ok := m.conns[pkg.UserID]
	if !ok {
		return
	}
	select {
	case agent.mq <- pkg.Body:
	default:
		log.Errorf("")
		agent.close()
	}
}
