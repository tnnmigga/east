package nats

import (
	"east/core/codec"
	"east/core/define"
	"east/core/iconf"
	"east/core/log"
	"east/core/message"
	"east/core/pb"
	"east/core/util"
	"time"
)

func (m *Module) initHandler() {
	message.RegisterHandler(m, m.onPackage)
	message.RegisterHandler(m, m.onBroadcastPackage)
	message.RegisterHandler(m, m.onRandomCastPackage)
	message.RegisterHandler(m, m.onRPCPequest)
}

func (m *Module) onPackage(pkg *define.Package) {
	b := codec.Encode(pkg.Body)
	netPkg := &pb.Package{
		Body: b,
	}
	_, err := m.js.PublishAsync(castSubject(pkg.ServerID), codec.Encode(netPkg))
	if err != nil {
		log.Errorf("nats publish error %v", err)
	}
}

func (m *Module) onBroadcastPackage(pkg *define.BroadcastPackage) {
	b := codec.Encode(pkg.Body)
	netPkg := &pb.Package{
		Body: b,
	}
	err := m.conn.Publish(broadcastSubject(pkg.ServerType), codec.Encode(netPkg))
	if err != nil {
		log.Errorf("nats publish error %v", err)
	}
}

func (m *Module) onRandomCastPackage(pkg *define.RandomCastPackage) {
	b := codec.Encode(pkg.Body)
	netPkg := &pb.Package{
		Body: b,
	}
	err := m.conn.Publish(randomCastSubject(pkg.ServerType), codec.Encode(netPkg))
	if err != nil {
		log.Errorf("nats publish error %v", err)
	}
}

func (m *Module) onRPCPequest(req *define.RPCRequest) {
	b := codec.Encode(req.Req)
	netPkg := &pb.Package{
		Body: b,
	}
	go util.ExecAndRecover(func() {
		msg, err := m.conn.Request(rpcSubject(req.ServerID), codec.Encode(netPkg), time.Duration(iconf.Int64("rpc-wait-time", 10))*time.Second)
		if err == nil {
			req.Resp, req.Err = codec.Decode(msg.Data)
		} else {
			req.Err = err
		}
		req.Module.MQ() <- req
	})
}
