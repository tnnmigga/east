package nats

import (
	"eden/core/codec"
	"eden/core/message"
	"eden/core/module"
	"eden/core/pb"
)

func (m *Module) initHandler() {
	module.RegisterHandler(m, m.onPackage)
}

func (m *Module) onPackage(pkg *message.Package) {
	b := codec.Encode(pkg.Body)
	netPkg := &pb.Package{
		Module: pkg.Module,
		Body:   b,
	}
	m.conn.Publish(defaultTopic(pkg.ServerID), codec.Encode(netPkg))
}
