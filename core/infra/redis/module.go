package redis

import (
	"east/core/basic"
	"east/core/conf"
	"east/core/idef"

	"github.com/go-redis/redis/v8"
)

type module struct {
	*basic.Module
	cli *redis.Client
}

func New() idef.IModule {
	m := &module{
		Module: basic.New("redis", basic.DefaultMQLen),
	}
	m.After(idef.ServerStateInit, m.afterInit)
	return m
}

func (m *module) afterInit() error {
	m.cli = redis.NewClient(&redis.Options{
		Addr: conf.String("redis.addr", "localhost:6379"),
	})
	return nil
}