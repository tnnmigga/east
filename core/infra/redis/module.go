package redis

import (
	"context"
	"east/core/basic"
	"east/core/conf"
	"east/core/idef"
	"time"

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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err := m.cli.Ping(ctx).Result()
	cancel()
	return err
}