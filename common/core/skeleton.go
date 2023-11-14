package core

import (
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

type CoreConfigs struct {
	IP        string
	Port      int32
	RedisOpts *redis.Options
	// MongoOpts *mongo.Co
}

type Skeleton struct {
	Redis redis.Cmdable
	Mongo *mongo.Client
}

func Init(configs *CoreConfigs) *Skeleton {
	skeleton := &Skeleton{
		Redis: redis.NewClient(configs.RedisOpts),
	}
	return skeleton
}
