package main

import (
	"context"
	"time"
	"github.com/redis/go-redis/v9"
)

type LeakyBucket struct {
	rdb      *redis.Client
	capacity int
	leakRate float64
}

func NewLeakyBucket(rdb *redis.Client, capacity int, leakRate float64) *LeakyBucket {
	return &LeakyBucket{
		rdb:      rdb,
		capacity: capacity,
		leakRate: leakRate,
	}
}

func (l *LeakyBucket) Allow(ctx context.Context, key string) (bool, error) {
	now := time.Now().UnixMilli()

	res, err := l.rdb.Eval(ctx, luaScript,
		[]string{key},
		l.capacity,
		l.leakRate,
		now,
	).Int()

	if err != nil {
		return false, err
	}
	return res == 1, nil
}

const luaScript = `
local data = redis.call("HMGET", KEYS[1], "level", "last_ts")
local level = tonumber(data[1]) or 0
local last_ts = tonumber(data[2]) or ARGV[3]

local now = tonumber(ARGV[3])
local leak_rate = tonumber(ARGV[2])
local capacity = tonumber(ARGV[1])

local elapsed = (now - last_ts) / 1000
local leaked = elapsed * leak_rate
level = math.max(0, level - leaked)
if level + 1 > capacity then
	return 0
end

level = level + 1

redis.call("HMSET", KEYS[1], "level", level, "last_ts", now)
redis.call("PEXPIRE", KEYS[1], math.ceil(capacity / leak_rate * 1000))

return 1
`