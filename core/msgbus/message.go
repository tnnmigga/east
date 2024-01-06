package msgbus

import (
	"east/core/conf"
	"east/core/idef"
	"east/core/log"
	"east/core/sys"
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
		rpcMaxWaitTime = time.Duration(conf.Int64("rpc-wait-time", 10)) * time.Second
	})
}

var (
	recvers        map[reflect.Type][]IRecver
	rpcMaxWaitTime time.Duration
)

func init() {
	recvers = map[reflect.Type][]IRecver{}
}

type IRecver interface {
	Name() string
	Assign(any)
}

func Cast(serverID uint32, msg any, opts ...castOpt) {
	if serverID == conf.ServerID() {
		dispatchMsg(msg, opts...)
		return
	}
	if nonuse, find := findCastOpt[bool](opts, keyNonuseStream); nonuse && find { // 不使用流
		Cast(conf.ServerID(), &idef.CastPackage{
			ServerID: serverID,
			Body:     msg,
		}, opts...)
		return
	}
	Cast(conf.ServerID(), &idef.StreamCastPackage{
		ServerID: serverID,
		Body:     msg,
	}, opts...)
}

func Broadcast(serverType string, msg any) {
	pkg := &idef.BroadcastPackage{
		ServerType: serverType,
		Body:       msg,
	}
	Cast(conf.ServerID(), pkg)
}

func RPC[T any](m idef.IModule, serverID uint32, req any, cb func(resp T, err error)) {
	if serverID == conf.ServerID() {
		localCall(m, req, warpCb(cb))
		return
	}
	rpcReq := &idef.RPCPackage{
		Module:   m,
		ServerID: serverID,
		Req:      req,
		Cb:       warpCb(cb),
	}
	Cast(conf.ServerID(), rpcReq)
}

func localCall(m idef.IModule, req any, cb func(resp any, err error)) {
	recvs, ok := recvers[reflect.TypeOf(req)]
	if !ok {
		log.Errorf("recvs not fuound %v", util.StructName(req))
		return
	}
	sys.Go(func() {
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
		timer := time.NewTimer(rpcMaxWaitTime)
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
			log.Panicf("async call resp type error, need %s, cur %s", util.StructName(new(T)), util.StructName(pkg))
		}
		cb(resp, err)
	}
}

func dispatchMsg(msg any, opts ...castOpt) {
	recvs, ok := recvers[reflect.TypeOf(msg)]
	if !ok {
		log.Errorf("message cast recv not fuound %v", util.StructName(msg))
		return
	}
	modName, oneOfMod := findCastOpt[string](opts, keyOneOfModules)
	for _, recv := range recvs {
		if oneOfMod && modName != recv.Name() {
			continue
		}
		recv.Assign(msg)
	}
}
