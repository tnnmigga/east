package nats

import (
	"context"
	"east/core/codec"
	"east/core/iconf"
	"east/core/idef"
	"east/core/log"
	"east/core/module"
	"east/core/msgbus"
	"east/core/sys"
	"east/core/util"
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
	js     jetstream.JetStream
	stream jetstream.Stream
	cons   jetstream.Consumer
}

func New(name string) idef.IModule {
	conn, err := nats.Connect(
		iconf.String("nats-url", nats.DefaultURL),
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
		Module: module.New(name, iconf.Int32("nats-mq-len", module.DefaultMQLen)),
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
	defer stop()
	m.Module.Run()
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

func randomCastSubject(serverType string) string {
	return fmt.Sprintf("randomcast.%s", serverType)
}

func rpcSubject(serverID uint32) string {
	return fmt.Sprintf("rpc.%d", serverID)
}

func (m *Module) initSubcribe() (stop func(), err error) {
	m.js, err = jetstream.New(m.conn)
	if err != nil {
		return nil, err
	}
	m.stream, err = m.js.Stream(context.Background(), castStreamName)
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
	queueSub, err := m.conn.QueueSubscribe(randomCastSubject(iconf.ServerType()), iconf.ServerType(), m.recv)
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
		queueSub.Drain()
		rpcSub.Drain()
		<-m.js.PublishAsyncComplete()
		m.conn.Close()
	}, nil
}

func (m *Module) streamRecv(msg jetstream.Msg) {
	defer util.RecoverPanic()
	msg.Ack()
	pkg, err := codec.Decode(msg.Data())
	if err != nil {
		log.Errorf("nats streamRecv decode msg error: %v", err)
		return
	}
	msgbus.Cast(iconf.ServerID(), pkg)
}

func (m *Module) recv(msg *nats.Msg) {
	defer util.RecoverPanic()
	pkg, err := codec.Decode(msg.Data)
	if err != nil {
		log.Errorf("nats recv decode msg error: %v", err)
		return
	}
	msgbus.Cast(iconf.ServerID(), pkg)
}

func (m *Module) rpc(msg *nats.Msg) {
	defer util.RecoverPanic()
	pkg, err := codec.Decode(msg.Data)
	if err != nil {
		log.Errorf("nats rpc decode msg error: %v", err)
		return
	}
	rpcMsg := &idef.RPCPackage{
		Req:  pkg,
		Resp: make(chan any, 1),
		Err:  make(chan error, 1),
	}
	msgbus.Cast(iconf.ServerID(), rpcMsg)
	sys.Go(func() {
		timer := time.After(time.Duration(iconf.Int64("rpc-wait-time", 10)) * time.Second)
		select {
		case <-timer:
			log.Errorf("nats rpc call timeout %v", util.StructName(rpcMsg.Req))
		case resp := <-rpcMsg.Resp:
			b := codec.Encode(resp)
			m.conn.Publish(msg.Reply, b)
		case err := <-rpcMsg.Err:
			log.Errorf("nats rpc call %v error %v", util.StructName(rpcMsg.Req), err)
		}
	})
}
