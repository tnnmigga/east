package nats

import (
	"east/core/codec"
	"east/core/conf"
	"east/core/idef"
	"east/core/log"
	"east/core/msgbus"
	"east/core/sys"
	"errors"
	"time"
)

func (com *component) initHandler() {
	msgbus.RegisterHandler(com, com.onCastPackage)
	msgbus.RegisterHandler(com, com.onStreamCastPackage)
	msgbus.RegisterHandler(com, com.onBroadcastPackage)
	msgbus.RegisterHandler(com, com.onRandomCastPackage)
	msgbus.RegisterHandler(com, com.onRPCPackage)
}

func (com *component) onCastPackage(pkg *idef.CastPackage) {
	b := codec.Encode(pkg.Body)
	err := com.conn.Publish(castSubject(pkg.ServerID), b)
	if err != nil {
		log.Errorf("nats publish error %v", err)
	}
}

func (com *component) onStreamCastPackage(pkg *idef.StreamCastPackage) {
	b := codec.Encode(pkg.Body)
	_, err := com.js.PublishAsync(streamCastSubject(pkg.ServerID), b)
	if err != nil {
		log.Errorf("nats publish error %v", err)
	}
}

func (com *component) onBroadcastPackage(pkg *idef.BroadcastPackage) {
	b := codec.Encode(pkg.Body)
	err := com.conn.Publish(broadcastSubject(pkg.ServerType), b)
	if err != nil {
		log.Errorf("nats publish error %v", err)
	}
}

func (com *component) onRandomCastPackage(pkg *idef.RandomCastPackage) {
	b := codec.Encode(pkg.Body)
	err := com.conn.Publish(randomCastSubject(pkg.ServerType), b)
	if err != nil {
		log.Errorf("nats publish error %v", err)
	}
}

func (com *component) onRPCPackage(req *idef.RPCPackage) {
	b := codec.Encode(req.Req)
	sys.Go(func() {
		resp := &idef.RPCResponse{
			Compt: req.Compt,
			Req:   req.Req,
			Cb:    req.Cb,
		}
		defer req.Compt.Assign(resp)
		msg, err := com.conn.Request(rpcSubject(req.ServerID), b, time.Duration(conf.Int64("rpc-wait-time", 10))*time.Second)
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
