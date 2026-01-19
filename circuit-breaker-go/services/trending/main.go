package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"circuit-breaker-go/pkg/circuitbreaker"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const serviceName = "trending"

type post struct {
	ID int64 `json:"id"`
	Email string `json:"email"`
	Content string `json:"content"`
	Views int `json:"views"`
	CreatedAt time.Time `json:"created_at"`
}

var (
	cache *circuitbreaker.StateCache
	postsClient *circuitbreaker.ProtectedHTTPClient
)

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func initCircuitBreakers() {
	cache = circuitbreaker.NewStateCache()

	redisConfig := circuitbreaker.DefaultRedisConfig()
	redisConfig.Addr = getEnv("REDIS_URL", "localhost:6379")
	pubsub := circuitbreaker.NewRedisPubSub(redisConfig)

	if err := pubsub.Ping(); err != nil {
		log.Println("[Trending] Warning: Redis connection failed: ", err)
	}

	pubsub.Subscribe(cache, serviceName)

	config := circuitbreaker.DefaultConfig()

	postsBreaker := circuitbreaker.NewCircuitBreaker(
		serviceName,
		"posts",
		config,
		cache,
		pubsub,
	)
	postsClient = circuitbreaker.NewProtectedHTTPClient(postsBreaker, config.RequestTimeout)

	log.Printf("[Trending] Circuit breaker initialized for posts service")
}

func main() {
	godotenv.Load()

	initCircuitBreakers()

	router := gin.Default()

	router.GET("/health/circuit-breakers", func(ctx *gin.Context) {
		states := cache.GetAll()
		result := make(map[string]interface{})
		for key, state := range states {
			result[string(key)] = gin.H{
				"state": state.State.String(),
				"failure_count": state.FailureCount,
				"last_change": state.LastStateChange,
			}
		}
		ctx.JSON(http.StatusOK, result)
	})

	router.GET("/trending", func(ctx *gin.Context) {
		POSTSAPI_URL := getEnv("POSTSAPI_URL", "http://localhost:7001")
		limit := ctx.Query("limit")
		if limit == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "limit not set in query param",
			})
			return
		}

		resp, err := postsClient.Get(POSTSAPI_URL + "/topPosts?limit=" + limit)
		if err != nil {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{
				"error": "posts service unavailable: " + err.Error(),
				"circuit": "may be open",
			})
			return
		}
		defer resp.Body.Close()

		var posts []post
		err1 := json.NewDecoder(resp.Body).Decode(&posts)
		if err1 != nil {
			fmt.Println("Error parsing json: ", err1)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err1.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, posts)
	})

	router.Run(":7002")
}
