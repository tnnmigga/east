package nats

import (
	"context"
	"east/core/codec"
	"east/core/iconf"
	"east/core/idef"
	"east/core/infra"
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
	conn         *nats.Conn
	js           jetstream.JetStream
	stream       jetstream.Stream
	cons         jetstream.Consumer
	consCtx      jetstream.ConsumeContext
	castSub      *nats.Subscription
	broadcastSub *nats.Subscription
	queueSub     *nats.Subscription
	rpcSub       *nats.Subscription
}

func New() idef.IModule {
	m := &Module{
		Module: module.New(infra.ModTypNats, iconf.Int32("nats-mq-len", module.DefaultMQLen)),
	}
	m.initHandler()
	m.After(idef.ServerStateInit, m.afterInit)
	m.After(idef.ServerStateRun, m.afterRun)
	m.Before(idef.ServerStateStop, m.beforeStop)
	m.After(idef.ServerStateStop, m.afterStop)
	return m
}

func (m *Module) afterInit() error {
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
		return err
	}
	m.conn = conn
	m.js, err = jetstream.New(m.conn)
	if err != nil {
		return err
	}
	m.stream, err = m.js.Stream(context.Background(), castStreamName)
	if err != nil {
		return err
	}
	return nil
}

func (m *Module) afterRun() (err error) {
	m.cons, err = m.stream.CreateOrUpdateConsumer(context.Background(), jetstream.ConsumerConfig{
		Durable:       fmt.Sprintf("%s-%d", iconf.ServerType(), iconf.ServerID()),
		FilterSubject: streamCastSubject(iconf.ServerID()),
	})
	if err != nil {
		return err
	}
	m.consCtx, err = m.cons.Consume(m.streamRecv)
	if err != nil {
		return err
	}
	m.castSub, err = m.conn.Subscribe(castSubject(iconf.ServerID()), m.recv)
	if err != nil {
		return err
	}
	m.broadcastSub, err = m.conn.Subscribe(broadcastSubject(iconf.ServerType()), m.recv)
	if err != nil {
		return err
	}
	m.queueSub, err = m.conn.QueueSubscribe(randomCastSubject(iconf.ServerType()), iconf.ServerType(), m.recv)
	if err != nil {
		return err
	}
	m.rpcSub, err = m.conn.Subscribe(rpcSubject(iconf.ServerID()), m.rpc)
	if err != nil {
		return err
	}
	return nil
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

func (m *Module) beforeStop() error {
	m.consCtx.Stop()
	m.castSub.Drain()
	m.broadcastSub.Drain()
	m.queueSub.Drain()
	m.rpcSub.Drain()
	return nil
}

func (m *Module) afterStop() error {
	<-m.js.PublishAsyncComplete()
	m.conn.Close()
	return nil
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
