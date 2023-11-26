package nats

import (
	"eden/core/codec"
	"eden/core/iconf"
	"eden/core/log"
	"eden/core/message"
	"eden/core/module"
	"eden/core/pb"
	"eden/core/util"
	"time"
)

func (m *Module) initHandler() {
	module.RegisterHandler(m.Module, m.onPackage)
	module.RegisterHandler(m.Module, m.onBroadcastPackage)
}

func (m *Module) onPackage(pkg *message.Package) {
	b := codec.Encode(pkg.Body)
	netPkg := &pb.Package{
		Module: pkg.Module,
		Body:   b,
	}
	err := m.conn.Publish(castTopic(pkg.ServerID), codec.Encode(netPkg))
	if err != nil {
		log.Errorf("nats publish error %v", err)
	}
}

func (m *Module) onBroadcastPackage(pkg *message.BroadcastPackage) {
	b := codec.Encode(pkg.Body)
	netPkg := &pb.Package{
		Module: pkg.Module,
		Body:   b,
	}
	err := m.conn.Publish(broadcastTopic(pkg.ServerType), codec.Encode(netPkg))
	if err != nil {
		log.Errorf("nats publish error %v", err)
	}
}

func (m *Module) onRPCPequest(req *message.RPCRequest) {
	b := codec.Encode(req.Req)
	netPkg := &pb.Package{
		Module: req.Module,
		Body:   b,
	}
	go util.ExecAndRecover(func() {
		msg, err := m.conn.Request(rpcTopic(req.ServerID), codec.Encode(netPkg), time.Duration(iconf.Int64("rpc-wait-time", 10))*time.Second)
		if err != nil {
			req.Resp, err = codec.Decode(msg.Data)
		} else {
			req.Err = err
		}
		message.Cast(iconf.ServerID(), req.Caller, req)
	})
}
