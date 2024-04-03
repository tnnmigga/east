package link

import (
	"east/core/basic"
	"east/core/codec"
	"east/core/conf"
	"east/core/idef"
	"east/core/log"
	"east/core/msgbus"
	"errors"
	fmt "fmt"

	"github.com/nats-io/nats.go"
)

func (m *module) initHandler() {
	msgbus.RegisterHandler(m, m.onCastPackage)
	msgbus.RegisterHandler(m, m.onStreamCastPackage)
	msgbus.RegisterHandler(m, m.onBroadcastPackage)
	msgbus.RegisterHandler(m, m.onRandomCastPackage)
	msgbus.RegisterHandler(m, m.onRPContext)
}

func (m *module) onCastPackage(pkg *idef.CastPackage) {
	b := codec.Encode(pkg.Body)
	err := m.conn.Publish(castSubject(pkg.ServerID), b)
	if err != nil {
		log.Errorf("onCastPackage error %v", err)
	}
}

func (m *module) onStreamCastPackage(pkg *idef.StreamCastPackage) {
	b := codec.Encode(pkg.Body)
	// _, err := m.js.PublishAsync(streamCastSubject(pkg.ServerID), b)
	msg := &nats.Msg{
		Subject: streamCastSubject(pkg.ServerID),
		Data:    b,
	}
	if len(pkg.Header) > 0 {
		msg.Header = nats.Header{}
		for key, value := range pkg.Header {
			msg.Header.Set(key, value)
		}
	}
	_, err := m.js.PublishMsgAsync(msg)
	if err != nil {
		log.Errorf("onStreamCastPackage error %v", err)
	}
}

func (m *module) onBroadcastPackage(pkg *idef.BroadcastPackage) {
	b := codec.Encode(pkg.Body)
	err := m.conn.Publish(broadcastSubject(pkg.ServerType), b)
	if err != nil {
		log.Errorf("onBroadcastPackage error %v", err)
	}
}

func (m *module) onRandomCastPackage(pkg *idef.RandomCastPackage) {
	b := codec.Encode(pkg.Body)
	err := m.conn.Publish(randomCastSubject(pkg.ServerType), b)
	if err != nil {
		log.Errorf("onRandomCastPackage error %v", err)
	}
}

func (m *module) onRPContext(ctx *idef.RPCContext) {
	b := codec.Encode(ctx.Req)
	basic.Go(func() {
		resp := &idef.RPCResponse{
			Module: ctx.Caller,
			Req:    ctx.Req,
			Cb:     ctx.Cb,
			Resp:   ctx.Resp,
		}
		defer ctx.Caller.Assign(resp)
		msg, err := m.conn.Request(rpcSubject(ctx.ServerID), b, conf.MaxRPCWaitTime)
		if err != nil {
			resp.Err = err
			return
		}
		data, err := codec.Decode(msg.Data)
		if err != nil {
			resp.Err = fmt.Errorf("RPCPkg decode error: %v", err)
			return
		}
		rpcResp := data.(*RPCResponse)
		if len(rpcResp.Err) != 0 {
			resp.Err = errors.New(rpcResp.Err)
			return
		}
		resp.Err = codec.Unmarshal(rpcResp.Data, resp.Resp)
	})
}
