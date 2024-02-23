package nats

import (
	"east/core/codec"
	"east/core/conf"
	"east/core/idef"
	"east/core/log"
	"east/core/msgbus"
	"east/core/sys"
	"errors"
)

func (m *module) initHandler() {
	msgbus.RegisterHandler(m, m.onCastPackage)
	msgbus.RegisterHandler(m, m.onStreamCastPackage)
	msgbus.RegisterHandler(m, m.onBroadcastPackage)
	msgbus.RegisterHandler(m, m.onRandomCastPackage)
	msgbus.RegisterHandler(m, m.onRPCPackage)
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

func (m *module) onRPCPackage(req *idef.RPCPackage) {
	b := codec.Encode(req.Req)
	sys.Go(func() {
		resp := &idef.RPCResponse{
			Module: req.Module,
			Req:    req.Req,
			Cb:     req.Cb,
			Resp:   req.Resp,
		}
		defer req.Module.Assign(resp)
		msg, err := m.conn.Request(rpcSubject(req.ServerID), b, conf.MaxRPCWaitTime)
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
		resp.Err = codec.Unmarshal(rpcResp.Data, resp.Resp)
	})
}
