package nats

import (
	"context"
	"east/core/codec"
	"east/core/com"
	"east/core/conf"
	"east/core/idef"
	"east/core/infra"
	"east/core/log"
	"east/core/msgbus"
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

type module struct {
	*com.Component
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
	m := &module{
		Component: com.New(infra.ModTypNats, conf.Int32("nats-mq-len", com.DefaultMQLen)),
	}
	codec.Register((*RPCResponse)(nil))
	m.initHandler()
	m.After(idef.ServerStateInit, m.afterInit)
	m.After(idef.ServerStateRun, m.afterRun)
	m.Before(idef.ServerStateStop, m.beforeStop)
	m.After(idef.ServerStateStop, m.afterStop)
	return m
}

func (m *module) afterInit() error {
	conn, err := nats.Connect(
		conf.String("nats-url", nats.DefaultURL),
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

func (m *module) afterRun() (err error) {
	m.cons, err = m.stream.CreateOrUpdateConsumer(context.Background(), jetstream.ConsumerConfig{
		Durable:       fmt.Sprintf("%s-%d", conf.ServerType(), conf.ServerID()),
		FilterSubject: streamCastSubject(conf.ServerID()),
	})
	if err != nil {
		return err
	}
	m.consCtx, err = m.cons.Consume(m.streamRecv)
	if err != nil {
		return err
	}
	m.castSub, err = m.conn.Subscribe(castSubject(conf.ServerID()), m.recv)
	if err != nil {
		return err
	}
	m.broadcastSub, err = m.conn.Subscribe(broadcastSubject(conf.ServerType()), m.recv)
	if err != nil {
		return err
	}
	m.queueSub, err = m.conn.QueueSubscribe(randomCastSubject(conf.ServerType()), conf.ServerType(), m.recv)
	if err != nil {
		return err
	}
	m.rpcSub, err = m.conn.Subscribe(rpcSubject(conf.ServerID()), m.rpc)
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

func (m *module) beforeStop() error {
	m.consCtx.Stop()
	m.castSub.Drain()
	m.broadcastSub.Drain()
	m.queueSub.Drain()
	m.rpcSub.Drain()
	return nil
}

func (m *module) afterStop() error {
	<-m.js.PublishAsyncComplete()
	m.conn.Close()
	return nil
}

func (m *module) streamRecv(msg jetstream.Msg) {
	defer util.RecoverPanic()
	msg.Ack()
	pkg, err := codec.Decode(msg.Data())
	if err != nil {
		log.Errorf("nats streamRecv decode msg error: %v", err)
		return
	}
	msgbus.Cast(conf.ServerID(), pkg)
}

func (m *module) recv(msg *nats.Msg) {
	defer util.RecoverPanic()
	pkg, err := codec.Decode(msg.Data)
	if err != nil {
		log.Errorf("nats recv decode msg error: %v", err)
		return
	}
	msgbus.Cast(conf.ServerID(), pkg)
}

func (m *module) rpc(msg *nats.Msg) {
	defer util.RecoverPanic()
	req, err := codec.Decode(msg.Data)
	rpcResp := &RPCResponse{}
	if err != nil {
		rpcResp.Err = fmt.Sprintf("req decode msg error: %v", err)
		m.conn.Publish(msg.Reply, codec.Encode(rpcResp))
		return
	}
	msgbus.RPC[any](m, conf.ServerID(), req, func(resp any, err error) {
		if err != nil {
			rpcResp.Err = err.Error()
		} else {
			rpcResp.Data = codec.Encode(resp)
		}
		m.conn.Publish(msg.Reply, codec.Encode(rpcResp))
	})
}
