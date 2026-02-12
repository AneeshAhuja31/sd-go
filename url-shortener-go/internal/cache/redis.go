package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)


func InitRedis(ctx context.Context, host string, port int)(*redis.Client,error){
	rdb := redis.NewClient(&redis.Options{
		Addr: host + ":" + fmt.Sprint(port),
		Password: "",
		DB: 0,
	})
	
	return rdb, nil
}

