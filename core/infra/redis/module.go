package redis

import (
	"east/core/module"

	"github.com/go-redis/redis"
)

type Module struct {
	*module.Module
	cli *redis.Client
}

func NewModule() {

}
