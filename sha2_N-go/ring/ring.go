package ring

import (
	"context"
	"sha-go/config"
	"sha-go/hash"
	"sha-go/node"
	"github.com/redis/go-redis/v9"
)

type Ring struct {
	RedisConn *redis.Client
	Ctx context.Context
	Ring []node.Node
}

func MakeRing() *Ring {
	redisClient,ctx := config.InitRedis("localhost", 6379)
	newRing := make([]node.Node,hash.TotalSlots)
	return &Ring{
		RedisConn: redisClient,
		Ctx: ctx,
		Ring: newRing,
	}
}
