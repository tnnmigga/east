package tcpagent

import (
	"east/core/log"
	"east/core/msgbus"
	"east/pb"
)

func (com *component) initHandler() {
	msgbus.RegisterHandler(com, com.onS2CPackage)
}

func (com *component) onS2CPackage(pkg *pb.S2CPackage) {
	com.RLock()
	defer com.RUnlock()
	agent, ok := com.conns[pkg.UserID]
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
