package ring

import (
	"sha-go/node"
	"sha-go/config"
	"github.com/redis/go-redis/v9"
	"context"
)

type Ring struct {
	RedisConn *redis.Client
	Ctx context.Context
	Ring []node.Node
}

func MakeRing()*Ring{
	redisClient,ctx := config.InitRedis("localhost",6379)
	newRing := make([]node.Node, 16)
	return &Ring{
		RedisConn: redisClient,
		Ctx: ctx,
		Ring: newRing,
	}
}