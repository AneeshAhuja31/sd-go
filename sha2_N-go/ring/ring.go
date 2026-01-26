package ring

import (
	"sha-go/node"
	"sha-go/config"
	"github.com/redis/go-redis/v9"
)

type Ring struct {
	redisConn *redis.Client
	ring []node.Node
}

func MakeRing()*ring{
	
}