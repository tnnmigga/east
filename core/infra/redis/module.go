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
	m.initHandler()
	m.cli = redis.NewClient(&redis.Options{
		Addr:     conf.String("redis.address", "localhost:6379"),
		Password: conf.String("redis.password", ""),
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := m.cli.Ping(ctx).Result()
	return err
}
