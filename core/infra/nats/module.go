package nats

import (
	"eden/core/codec"
	"eden/core/iconf"
	"eden/core/log"
	"eden/core/message"
	"eden/core/module"
	"eden/core/pb"
	"eden/core/util"
	"fmt"
	"reflect"
	"time"

	"github.com/nats-io/nats.go"
)

const ModuleName = "nats"

type Module struct {
	*module.Module
	conn *nats.Conn
}

var (
	netPackageType = reflect.TypeOf((*pb.Package)(nil))
)

func NewModule(name string) module.IModule {
	conn, err := nats.Connect(iconf.String("nats-url"))
	if err != nil {
		panic(err)
	}
	m := &Module{
		Module: module.NewModule(name, iconf.Int32("mq-len")),
		conn:   conn,
	}
	m.initHandler()
	return m
}

func castTopic(serverID uint32) string {
	return fmt.Sprintf("cast.%s.%d", iconf.ServerType(), serverID)
}

func broadcastTopic(serverType string) string {
	return fmt.Sprintf("broadcast.%s", serverType)
}

func rpcTopic(serverID uint32) string {
	return fmt.Sprintf("rpc.%s.%d", iconf.ServerType(), serverID)
}

func (m *Module) Run() {
	m.conn.Subscribe(castTopic(iconf.ServerID()), m.recv)
	m.conn.Subscribe(broadcastTopic(iconf.ServerType()), m.recv)
	m.conn.Subscribe(rpcTopic(iconf.ServerID()), m.rpc)
	m.Module.Run()
}

func (m *Module) recv(msg *nats.Msg) {
	defer util.RecoverPanic()
	pkg, err := unpack(msg.Data)
	if err != nil {
		log.Errorf("nats recv decode msg error: %v", err)
		return
	}
	message.Cast(pkg.ServerID, pkg.Module, pkg.Body)
}

func (m *Module) rpc(msg *nats.Msg) {
	defer util.RecoverPanic()
	pkg, err := unpack(msg.Data)
	if err != nil {
		log.Errorf("nats rpc decode msg error: %v", err)
		return
	}
	rpcMsg := &message.RPCPackage{
		Req:  pkg.Body,
		Resp: make(chan any, 1),
	}
	message.Cast(pkg.ServerID, pkg.Module, rpcMsg)
	go util.ExecAndRecover(func() {
		timer := time.After(time.Duration(iconf.Int64("rpc-wait-time", 10)) * time.Second)
		select {
		case <-timer:
			log.Errorf("nats rpc call timeout %v", util.ReflectName(rpcMsg.Req))
		case resp := <-rpcMsg.Resp:
			b := codec.Encode(resp)
			m.conn.Publish(msg.Reply, b)
		}
	})
}

func unpack(b []byte) (*message.Package, error) {
	pkg0, err := codec.Decode(b)
	if err != nil {
		return nil, fmt.Errorf("nats decode msg error: %v", err)
	}
	pkg := pkg0.(*pb.Package)
	body, err := codec.Decode(pkg.Body)
	if err != nil {
		return nil, fmt.Errorf("nats recv decode pkg error: %v", err)
	}
	return &message.Package{
		ServerID: iconf.ServerID(),
		Module:   pkg.Module,
		Body:     body,
	}, nil
}
