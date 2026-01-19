package circuitbreaker

import (
	"time"
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	FailureThreshold int
	SuccessThreshold int
	Timeout time.Duration
	HalfOpenMaxRequests int
	RequestTimeout time.Duration 
}

func DefaultConfig() Config{
	return Config{
		FailureThreshold: 5,
		SuccessThreshold: 3,
		Timeout: 30 * time.Second,
		HalfOpenMaxRequests: 1,
		RequestTimeout: 10 * time.Second,
	}
} 

type RedisConfig struct{
	Addr string
	Password string
	DB int
	Channel string
}

func DefaultRedisConfig() RedisConfig {
	godotenv.Load()
	redisAddr,ok := os.LookupEnv("REDIS_URL")
	if !ok {
		redisAddr = "localhost:6379"
	}
	redisPassword,ok := os.LookupEnv("REDIS_PASSWORD")
	if !ok{
		redisPassword = ""
	}
	return RedisConfig{
		Addr: redisAddr,
		Password: redisPassword,
		DB: 0,
		Channel: "circuit_breaker_state",
	}
}