package hashring

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type HashRing struct {
	Rdb *redis.Client
	Ctx context.Context
	NumWorkers int
}

func InitRedis() (*redis.Client, context.Context) {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URL"),
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal("Redis connection failed: ", err)
	}
	return rdb, ctx
}

func hashKey(key string) uint64 {
	h := sha256.Sum256([]byte(key))
	return binary.BigEndian.Uint64(h[:8])
}

func NewHashRing(rdb *redis.Client, ctx context.Context, numWorkers int, vnodesPerWorker int) *HashRing {
	rdb.Del(ctx, "hashring")

	for i := 0; i < numWorkers; i++ {
		for j := 0; j < vnodesPerWorker; j++ {
			vnodeKey := fmt.Sprintf("worker-%d-vnode-%d", i, j)
			score := hashKey(vnodeKey)
			rdb.ZAdd(ctx, "hashring", redis.Z{
				Score: float64(score),
				Member: fmt.Sprint(i),
			})
		}
	}
	log.Printf("Hash ring initialized with %d workers x %d vnodes = %d entries",
		numWorkers, vnodesPerWorker, numWorkers*vnodesPerWorker)
	return &HashRing{Rdb: rdb, Ctx: ctx, NumWorkers: numWorkers}
}

func (hr *HashRing) GetWorker(domain string) int {
	domainHash := hashKey(domain)

	results, err := hr.Rdb.ZRangeArgs(hr.Ctx, redis.ZRangeArgs{
		Key: "hashring",
		Start: fmt.Sprint(domainHash),
		Stop: "+inf",
		ByScore: true,
		Count: 1,
	}).Result()

	if err != nil || len(results) == 0 {
		results, _ = hr.Rdb.ZRangeArgs(hr.Ctx, redis.ZRangeArgs{
			Key: "hashring",
			Start: 0,
			Stop: 0,
		}).Result()
	}

	workerID, _ := strconv.Atoi(results[0])
	return workerID
}
