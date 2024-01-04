package nats

import (
	"east/core/codec"
	"east/core/iconf"
	"east/core/idef"
	"east/core/log"
	"east/core/msgbus"
	"east/core/sys"
	"errors"
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
		resp := &idef.AsyncCallResponse{
			Module: req.Module,
			Req:    req.Req,
			Cb:     req.Cb,
		}
		defer req.Module.Assign(resp)
		msg, err := m.conn.Request(rpcSubject(req.ServerID), b, time.Duration(iconf.Int64("rpc-wait-time", 10))*time.Second)
		if err != nil {
			resp.Err = err
			return
		}
		rpcResp0, err := codec.Decode(msg.Data)
		if err != nil {
			resp.Err = errors.New("RPCPkg decode error")
			return
		}
		rpcResp := rpcResp0.(*RPCResponse)
		if len(rpcResp.Err) != 0 {
			resp.Err = errors.New(rpcResp.Err)
			return
		}
		resp.Resp, resp.Err = codec.Decode(msg.Data)
	})
}
