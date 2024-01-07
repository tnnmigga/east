package nats

import (
	"context"
	"east/core/codec"
	"east/core/compt"
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

type component struct {
	*compt.Component
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

func New() idef.IComponent {
	com := &component{
		Component: compt.New(infra.ModTypNats, conf.Int32("nats-mq-len", compt.DefaultMQLen)),
	}
	codec.Register((*RPCResponse)(nil))
	com.initHandler()
	com.After(idef.ServerStateInit, com.afterInit)
	com.After(idef.ServerStateRun, com.afterRun)
	com.Before(idef.ServerStateStop, com.beforeStop)
	com.After(idef.ServerStateStop, com.afterStop)
	return com
}

func (com *component) afterInit() error {
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
	com.conn = conn
	com.js, err = jetstream.New(com.conn)
	if err != nil {
		return err
	}
	com.stream, err = com.js.Stream(context.Background(), castStreamName)
	if err != nil {
		return err
	}
	return nil
}

func (com *component) afterRun() (err error) {
	com.cons, err = com.stream.CreateOrUpdateConsumer(context.Background(), jetstream.ConsumerConfig{
		Durable:       fmt.Sprintf("%s-%d", conf.ServerType(), conf.ServerID()),
		FilterSubject: streamCastSubject(conf.ServerID()),
	})
	if err != nil {
		return err
	}
	com.consCtx, err = com.cons.Consume(com.streamRecv)
	if err != nil {
		return err
	}
	com.castSub, err = com.conn.Subscribe(castSubject(conf.ServerID()), com.recv)
	if err != nil {
		return err
	}
	com.broadcastSub, err = com.conn.Subscribe(broadcastSubject(conf.ServerType()), com.recv)
	if err != nil {
		return err
	}
	com.queueSub, err = com.conn.QueueSubscribe(randomCastSubject(conf.ServerType()), conf.ServerType(), com.recv)
	if err != nil {
		return err
	}
	com.rpcSub, err = com.conn.Subscribe(rpcSubject(conf.ServerID()), com.rpc)
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

func (com *component) beforeStop() error {
	com.consCtx.Stop()
	com.castSub.Drain()
	com.broadcastSub.Drain()
	com.queueSub.Drain()
	com.rpcSub.Drain()
	return nil
}

func (com *component) afterStop() error {
	<-com.js.PublishAsyncComplete()
	com.conn.Close()
	return nil
}

func (com *component) streamRecv(msg jetstream.Msg) {
	defer util.RecoverPanic()
	msg.Ack()
	pkg, err := codec.Decode(msg.Data())
	if err != nil {
		log.Errorf("nats streamRecv decode msg error: %v", err)
		return
	}
	msgbus.CastLocal(pkg)
}

func (com *component) recv(msg *nats.Msg) {
	defer util.RecoverPanic()
	pkg, err := codec.Decode(msg.Data)
	if err != nil {
		log.Errorf("nats recv decode msg error: %v", err)
		return
	}
	msgbus.CastLocal(pkg)
}

func (com *component) rpc(msg *nats.Msg) {
	defer util.RecoverPanic()
	req, err := codec.Decode(msg.Data)
	rpcResp := &RPCResponse{}
	if err != nil {
		rpcResp.Err = fmt.Sprintf("req decode msg error: %v", err)
		com.conn.Publish(msg.Reply, codec.Encode(rpcResp))
		return
	}
	msgbus.RPC[any](com, conf.ServerID(), req, func(resp any, err error) {
		if err != nil {
			rpcResp.Err = err.Error()
		} else {
			rpcResp.Data = codec.Encode(resp)
		}
		com.conn.Publish(msg.Reply, codec.Encode(rpcResp))
	})
}
