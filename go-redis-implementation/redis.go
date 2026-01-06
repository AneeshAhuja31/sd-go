package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func connectRedis(host string, port int) (*redis.Client,context.Context){
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr: host + ":" + fmt.Sprint(port),
		Password: "",
		DB: 0,
	})
	return rdb,ctx
}

func setDataRedis(users []byte,rdb *redis.Client,ctx context.Context){
	err := rdb.Set(ctx,"users",users,0).Err()
	if err != nil {
		panic(err)
	}
}

func fetchDataRedis(rdb *redis.Client,ctx context.Context)([]byte,error){
	value,err := rdb.Get(ctx,"users").Result()
	// var cachedUsers []User
	// json.Unmarshal([]byte(value),&cachedUsers)
	// fmt.Println(cachedUsers)
	return []byte(value),err
}