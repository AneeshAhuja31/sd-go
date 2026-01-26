package ring

import (
	"github.com/AneeshAhuja31/sha2^N-go/node"
	"github.com/redis/redis-go/v9"
)

type Ring struct {
	redisConn *redis.Client
}