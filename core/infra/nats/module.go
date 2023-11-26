package nats

import (
	"eden/core/codec"
	"eden/core/configs"
	"eden/core/log"
	"eden/core/message"
	"eden/core/module"
	"eden/core/pb"
	"fmt"
	"reflect"

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

func NewModule() module.IModule {
	conn, err := nats.Connect(configs.String("nats-url"))
	if err != nil {
		panic(err)
	}
	m := &Module{
		Module: module.NewModule("nats", configs.Int32("mq-len")),
		conn:   conn,
	}
	m.initHandler()
	return m
}

func (m *Module) recv(msg *nats.Msg) {
	netPkg, err := codec.Decode(msg.Data)
	if err != nil {
		log.Errorf("nats decode msg error: %v", err)
		return
	}
	switch reflect.TypeOf(netPkg) {
	case netPackageType:
		pkg0 := netPkg.(*pb.Package)
		body, err := codec.Decode(pkg0.Body)
		if err != nil {
			log.Errorf("nats decode pkg error: %v", err)
			return
		}
		message.Cast(configs.ServerID(), pkg0.Module, body)
	default:
		log.Errorf("invalid net package type")
	}
}

func defaultTopic(serverID uint32) string {
	return fmt.Sprintf("single.%s.%d", configs.ServerType(), serverID)
}

func broadcastTopic() string {
	return fmt.Sprintf("broadcast.%s", configs.ServerType())
}

func (m *Module) Run() {
	m.conn.Subscribe(defaultTopic(configs.ServerID()), m.recv)
	m.conn.Subscribe(broadcastTopic(), m.recv)
	m.Module.Run()
}
