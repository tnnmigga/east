package mredis

import (
	"eden/core/module"

	"github.com/go-redis/redis"
)

type Module struct {
	*module.Module
	cli *redis.Client
}

func NewModule() {

}
