package config

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func initRedis(host string,port int) (*redis.Client,context.Context){
	ctx := context.Background()
	redisClient := redis.NewClient(&redis.Options{
		Addr: host + ":" + fmt.Sprint(port),
		Password: "",
		DB: 0,
	})
	log.Println("Setup redis connection")
	return redisClient,ctx
}