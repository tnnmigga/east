package link

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/tnnmigga/nett/basic"
	"github.com/tnnmigga/nett/codec"
	"github.com/tnnmigga/nett/conf"
	"github.com/tnnmigga/nett/idef"
	"github.com/tnnmigga/nett/infra"
	"github.com/tnnmigga/nett/msgbus"
	"github.com/tnnmigga/nett/util"
	"github.com/tnnmigga/nett/zlog"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	castStreamName = "stream-cast"
)

type module struct {
	*basic.Module
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
		Module: basic.New(infra.ModNameLink, conf.Int32("nats.mq-len", basic.DefaultMQLen)),
	}
	codec.Register((*RPCResult)(nil))
	m.initHandler()
	m.After(idef.ServerStateInit, m.afterInit)
	m.After(idef.ServerStateRun, m.afterRun)
	m.Before(idef.ServerStateStop, m.beforeStop)
	m.After(idef.ServerStateStop, m.afterStop)
	return m
}

func (m *module) afterInit() error {
	conn, err := nats.Connect(
		conf.String("nats.url", nats.DefaultURL),
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(10),
		nats.ReconnectWait(time.Second),
		nats.ReconnectHandler(func(_ *nats.Conn) {
			zlog.Errorf("nats retry connect")
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
	zlog.Debug(msg.Headers())
	if expires := msg.Headers().Get(idef.ConstKeyExpires); expires != "" {
		// 检测部分不重要但有一定时效性的消息是否超时
		// 比如往客户端推送的实时消息
		// 超时后直接丢弃
		n, err := strconv.Atoi(expires)
		if err == nil && util.NowNs() > time.Duration(n) {
			zlog.Debugf("message expired")
			return
		}
	}
	pkg, err := codec.Decode(msg.Data())
	if err != nil {
		zlog.Errorf("nats streamRecv decode msg error: %v", err)
		return
	}
	msgbus.CastLocal(pkg)
}

func (m *module) recv(msg *nats.Msg) {
	defer util.RecoverPanic()
	pkg, err := codec.Decode(msg.Data)
	if err != nil {
		zlog.Errorf("nats recv decode msg error: %v", err)
		return
	}
	msgbus.CastLocal(pkg)
}

func (m *module) rpc(msg *nats.Msg) {
	defer util.RecoverPanic()
	req, err := codec.Decode(msg.Data)
	rpcResp := &RPCResult{}
	if err != nil {
		rpcResp.Err = fmt.Sprintf("req decode msg error: %v", err)
		m.conn.Publish(msg.Reply, codec.Encode(rpcResp))
		return
	}
	msgbus.RPC[any](m, conf.ServerID(), req, func(resp any, err error) {
		if err != nil {
			rpcResp.Err = err.Error()
		} else {
			rpcResp.Data = codec.Marshal(resp)
		}
		m.conn.Publish(msg.Reply, codec.Encode(rpcResp))
	})
}