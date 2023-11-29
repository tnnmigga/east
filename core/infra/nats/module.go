package nats

import (
	"context"
	"eden/core/codec"
	"eden/core/iconf"
	"eden/core/log"
	"eden/core/message"
	"eden/core/module"
	"eden/core/pb"
	"eden/core/util"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	ModuleName     = "nats"
	castStreamName = "stream-cast"
)

type Module struct {
	*module.Module
	conn   *nats.Conn
	stream jetstream.Stream
	cons   jetstream.Consumer
}

func New(name string) module.IModule {
	conn, err := nats.Connect(
		iconf.String("nats-url"),
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(10),
		nats.ReconnectWait(time.Second),
		nats.ReconnectHandler(func(_ *nats.Conn) {
			log.Errorf("nats retry connect")
		}),
	)
	if err != nil {
		panic(err)
	}
	m := &Module{
		Module: module.New(name, iconf.Int32("nats-mq-len", 100000)),
		conn:   conn,
	}
	m.initHandler()
	return m
}

func (m *Module) Run() {
	stop, err := m.initSubcribe()
	if err != nil {
		panic(err)
	}
	m.Module.Run()
	stop()
}

func castSubject(serverID uint32) string {
	return fmt.Sprintf("cast.%d", serverID)
}

func streamCastSubject(serverID uint32) string {
	return fmt.Sprintf("stream.cast.%d", serverID)
}

func broadcastSubject(serverType string) string {
	return fmt.Sprintf("broadcast.%s", serverType)
}

func rpcSubject(serverID uint32) string {
	return fmt.Sprintf("rpc.%d", serverID)
}

func (m *Module) initSubcribe() (stop func(), err error) {
	js, err := jetstream.New(m.conn)
	if err != nil {
		return nil, err
	}
	m.stream, err = js.Stream(context.Background(), castStreamName)
	if err != nil {
		return nil, err
	}
	m.cons, err = m.stream.CreateOrUpdateConsumer(context.Background(), jetstream.ConsumerConfig{
		Durable:       fmt.Sprintf("%s-%d", iconf.ServerType(), iconf.ServerID()),
		FilterSubject: streamCastSubject(iconf.ServerID()),
	})
	if err != nil {
		return nil, err
	}
	consCtx, err := m.cons.Consume(m.streamRecv)
	if err != nil {
		return nil, err
	}
	castSub, err := m.conn.Subscribe(castSubject(iconf.ServerID()), m.recv)
	if err != nil {
		return nil, err
	}
	broadcastSub, err := m.conn.Subscribe(broadcastSubject(iconf.ServerType()), m.recv)
	if err != nil {
		return nil, err
	}
	rpcSub, err := m.conn.Subscribe(rpcSubject(iconf.ServerID()), m.rpc)
	if err != nil {
		return nil, err
	}
	return func() {
		consCtx.Stop()
		castSub.Drain()
		broadcastSub.Drain()
		rpcSub.Drain()
	}, nil
}

func (m *Module) Close() {
	m.conn.Close()
	m.Module.Close()
}

func (m *Module) streamRecv(msg jetstream.Msg) {
	defer util.RecoverPanic()
	msg.Ack()
	pkg, err := unpack(msg.Data())
	if err != nil {
		log.Errorf("nats streamRecv decode msg error: %v", err)
		return
	}
	message.Cast(pkg.ServerID, pkg.Module, pkg.Body)
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
		Err:  make(chan error, 1),
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
		case err := <-rpcMsg.Err:
			log.Errorf("nats rpc call %v error %v", util.ReflectName(rpcMsg.Req), err)
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
