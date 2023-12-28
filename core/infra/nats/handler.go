package nats

import (
	"east/core/codec"
	"east/core/iconf"
	"east/core/idef"
	"east/core/log"
	"east/core/msgbus"
	"east/core/sys"
	"time"
)

func (m *module) initHandler() {
	msgbus.RegisterHandler(m, m.onCastPackage)
	msgbus.RegisterHandler(m, m.onStreamCastPackage)
	msgbus.RegisterHandler(m, m.onBroadcastPackage)
	msgbus.RegisterHandler(m, m.onRandomCastPackage)
	msgbus.RegisterHandler(m, m.onRPCRequest)
}

func (m *module) onCastPackage(pkg *idef.CastPackage) {
	b := codec.Encode(pkg.Body)
	err := m.conn.Publish(castSubject(pkg.ServerID), b)
	if err != nil {
		log.Errorf("nats publish error %v", err)
	}
}

func (m *module) onStreamCastPackage(pkg *idef.StreamCastPackage) {
	b := codec.Encode(pkg.Body)
	_, err := m.js.PublishAsync(streamCastSubject(pkg.ServerID), b)
	if err != nil {
		log.Errorf("nats publish error %v", err)
	}
}

func (m *module) onBroadcastPackage(pkg *idef.BroadcastPackage) {
	b := codec.Encode(pkg.Body)
	err := m.conn.Publish(broadcastSubject(pkg.ServerType), b)
	if err != nil {
		log.Errorf("nats publish error %v", err)
	}
}

func (m *module) onRandomCastPackage(pkg *idef.RandomCastPackage) {
	b := codec.Encode(pkg.Body)
	err := m.conn.Publish(randomCastSubject(pkg.ServerType), b)
	if err != nil {
		log.Errorf("nats publish error %v", err)
	}
}

func (m *module) onRPCRequest(req *idef.RPCRequest) {
	b := codec.Encode(req.Req)
	sys.Go(func() {
		msg, err := m.conn.Request(rpcSubject(req.ServerID), b, time.Duration(iconf.Int64("rpc-wait-time", 10))*time.Second)
		if err == nil {
			req.Resp, req.Err = codec.Decode(msg.Data)
		} else {
			req.Err = err
		}
		req.Module.MQ() <- req
	})
}
