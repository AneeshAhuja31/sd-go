package config

import (
	"github.com/redis/go-redis/v9"
	"fmt"
	"context"
)

func InitRedisClient(host string, port int)(*redis.Client,context.Context){
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr: host + ":" + fmt.Sprint(port),
		Password: "",
		DB: 0,
	})
	return client,ctx
}