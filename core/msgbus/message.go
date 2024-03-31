package msgbus

import (
	"east/core/basic"
	"east/core/conf"
	"east/core/idef"
	"east/core/log"
	"east/core/util"
	"errors"
	"time"

	"reflect"
)

var (
	ErrRPCTimeout = errors.New("rpc timeout")
)

func init() {
	conf.RegInitFn(func() {
		// rpcMaxWaitTime = time.Duration(conf.Int64("rpc-wait-time", 10)) * time.Second
	})
}

var (
	recvers map[reflect.Type][]IRecver
	// rpcMaxWaitTime time.Duration
)

func init() {
	recvers = map[reflect.Type][]IRecver{}
}

type IRecver interface {
	Name() string
	Assign(any)
}

func Cast(msg any, opts ...castOpt) {
	serverID, find := findCastOpt[uint32](opts, keyServerID)
	if !find || serverID == conf.ServerID() {
		CastLocal(msg, opts...)
	}
	if nonuse, find := findCastOpt[bool](opts, keyNonuseStream); nonuse && find { // 不使用流
		CastLocal(&idef.CastPackage{
			ServerID: serverID,
			Body:     msg,
		}, opts...)
		return
	}
	CastLocal(&idef.StreamCastPackage{
		ServerID: serverID,
		Body:     msg,
	}, opts...)
}

func CastLocal(msg any, opts ...castOpt) {
	recvs, ok := recvers[reflect.TypeOf(msg)]
	if !ok {
		log.Errorf("message cast recv not fuound %v", util.StructName(msg))
		return
	}
	modName, oneOfMod := findCastOpt[string](opts, keyOneOfMods)
	for _, recv := range recvs {
		if oneOfMod && modName != recv.Name() {
			continue
		}
		recv.Assign(msg)
	}
}

func Broadcast(serverType string, msg any) {
	pkg := &idef.BroadcastPackage{
		ServerType: serverType,
		Body:       msg,
	}
	CastLocal(pkg)
}

func RPC[T any](m idef.IModule, serverID uint32, req any, cb func(resp T, err error)) {
	if serverID == conf.ServerID() {
		localCall(m, req, warpCb(cb))
		return
	}
	rpcCtx := &idef.RPCContext{
		Caller:   m,
		ServerID: serverID,
		Req:      req,
		Resp:     util.New[T](),
		Cb:       warpCb(cb),
	}
	CastLocal(rpcCtx)
}

func localCall(m idef.IModule, req any, cb func(resp any, err error)) {
	recvs, ok := recvers[reflect.TypeOf(req)]
	if !ok {
		log.Errorf("recvs not fuound %v", util.StructName(req))
		return
	}
	basic.Go(func() {
		callReq := &idef.RPCRequest{
			Req:  req,
			Resp: make(chan any, 1),
			Err:  make(chan error, 1),
		}
		callResp := &idef.RPCResponse{
			Module: m,
			Req:    req,
			Cb:     warpCb(cb),
		}
		recvs[0].Assign(callReq)
		timer := time.NewTimer(conf.MaxRPCWaitTime)
		defer timer.Stop()
		select {
		case <-timer.C:
			callResp.Err = ErrRPCTimeout
		case callResp.Resp = <-callReq.Resp:
		case callResp.Err = <-callReq.Err:
		}
		m.Assign(callResp)
	})
}

func warpCb[T any](cb func(T, error)) func(any, error) {
	return func(pkg any, err error) {
		if err != nil {
			var empty T
			cb(empty, err)
			return
		}
		resp, ok := pkg.(T)
		if !ok {
			log.Panicf("rpc resp type error, need %s, cur %s", util.StructName(new(T)), util.StructName(pkg))
		}
		cb(resp, err)
	}
}
