package nats

import (
	"east/core/codec"
	"east/core/iconf"
	"east/core/idef"
	"east/core/log"
	"east/core/message"
	"east/core/sys"
	"time"
)

func (m *Module) initHandler() {
	message.RegisterHandler(m, m.onPackage)
	message.RegisterHandler(m, m.onBroadcastPackage)
	message.RegisterHandler(m, m.onRandomCastPackage)
	message.RegisterHandler(m, m.onRPCPequest)
}

func (m *Module) onPackage(pkg *idef.CastPackage) {
	b := codec.Encode(pkg.Body)
	_, err := m.js.PublishAsync(castSubject(pkg.ServerID), b)
	if err != nil {
		log.Errorf("nats publish error %v", err)
	}
}

func (m *Module) onBroadcastPackage(pkg *idef.BroadcastPackage) {
	b := codec.Encode(pkg.Body)
	err := m.conn.Publish(broadcastSubject(pkg.ServerType), b)
	if err != nil {
		log.Errorf("nats publish error %v", err)
	}
}

func (m *Module) onRandomCastPackage(pkg *idef.RandomCastPackage) {
	b := codec.Encode(pkg.Body)
	err := m.conn.Publish(randomCastSubject(pkg.ServerType), b)
	if err != nil {
		log.Errorf("nats publish error %v", err)
	}
}

func (m *Module) onRPCPequest(req *idef.RPCRequest) {
	b := codec.Encode(req.Req)
	sys.Go[sys.Call](func() {
		msg, err := m.conn.Request(rpcSubject(req.ServerID), b, time.Duration(iconf.Int64("rpc-wait-time", 10))*time.Second)
		if err == nil {
			req.Resp, req.Err = codec.Decode(msg.Data)
		} else {
			req.Err = err
		}
		req.Module.MQ() <- req
	})
}
