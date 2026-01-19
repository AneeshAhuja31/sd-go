package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
	"circuit-breaker-go/pkg/circuitbreaker"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const serviceName = "feed"

type post struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Content   string    `json:"content"`
	Views     int       `json:"views"`
	CreatedAt time.Time `json:"created_at"`
}

var cache *circuitbreaker.StateCache
var trendingClient *circuitbreaker.ProtectedHTTPClient
var recommendationsClient *circuitbreaker.ProtectedHTTPClient



func initCircuitBreakers() {
	cache = circuitbreaker.NewStateCache()
	redisConfig := circuitbreaker.DefaultRedisConfig()
	if redis_url, ok := os.LookupEnv("REDIS_URL"); ok {
		redisConfig.Addr = redis_url
	}
	pubsub := circuitbreaker.NewRedisPubSub(redisConfig)

	if err := pubsub.Ping(); err != nil {
		log.Println("[Feed] Warning: Redis connection failed: ", err)
	}

	pubsub.Subscribe(cache, serviceName)

	config := circuitbreaker.DefaultConfig()

	trendingBreaker := circuitbreaker.NewCircuitBreaker( //cricuit breaker with trending sserive
		serviceName,
		"trending",
		config,
		cache,
		pubsub,
	)
	trendingClient = circuitbreaker.NewProtectedHTTPClient(trendingBreaker, config.RequestTimeout)

	recommendationsBreaker := circuitbreaker.NewCircuitBreaker(
		serviceName,
		"recommendations", //circuit breaker with recommendations service
		config,
		cache,
		pubsub,
	)
	recommendationsClient = circuitbreaker.NewProtectedHTTPClient(recommendationsBreaker, config.RequestTimeout)

	log.Printf("[Feed] Circuit breakers initialized for trending and recommendations services")
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

	router.GET("/feed", func(ctx *gin.Context) {
		email := ctx.Query("email")
		if email == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "email not set in query param",
			})
			return
		}

		TRENDING_URL,ok := os.LookupEnv("TRENDING_URL")
		if !ok {
			TRENDING_URL = "http://localhost:7002"
		}
		RECOMMENDATIONS_URL,ok := os.LookupEnv("RECOMMENDATIONS_URL")
		if !ok {
			RECOMMENDATIONS_URL = "http://localhost:7003"
		}

		trendingResp, err := trendingClient.Get(TRENDING_URL + "/trending?limit=10")
		if err != nil {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{
				"error":   "failed to fetch trending posts: " + err.Error(),
				"circuit": "may be open",
			})
			return
		}
		defer trendingResp.Body.Close()

		var trendingPosts []post
		err = json.NewDecoder(trendingResp.Body).Decode(&trendingPosts)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to parse trending posts: " + err.Error(),
			})
			return
		}

		recommendedResp, err := recommendationsClient.Get(RECOMMENDATIONS_URL + "/recommendations?email=" + email + "&limit=20")
		if err != nil {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{
				"error":"failed to fetch recommendations: " + err.Error(),
				"circuit":"may be open",
			})
			return
		}
		defer recommendedResp.Body.Close()

		var recommendedPosts []post
		if err := json.NewDecoder(recommendedResp.Body).Decode(&recommendedPosts); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to parse recommendations: " + err.Error(),
			})
			return
		}

		seenIDs := make(map[int]bool)
		feedPosts := []post{}

		for _, p := range trendingPosts {
			if !seenIDs[p.ID] {
				seenIDs[p.ID] = true
				feedPosts = append(feedPosts, p)
			}
		}

		for _, p := range recommendedPosts {
			if !seenIDs[p.ID] {
				seenIDs[p.ID] = true
				feedPosts = append(feedPosts, p)
			}
		}

		ctx.JSON(http.StatusOK, feedPosts)
	})

	router.Run(":7004")
}
