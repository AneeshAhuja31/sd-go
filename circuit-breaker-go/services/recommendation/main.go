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

const serviceName = "recommendations"

type profile struct {
	ID int64 `json:"id"`
	Email string `json:"email"`
	Username  string    `json:"username"`
	DOB time.Time`json:"dob"`
	Bio string `json:"bio"`
	Hobbies []string `json:"hobbies"`
	CreatedAt time.Time `json:"created_at"`
}

type post struct {
	ID int64 `json:"id"`
	Email string `json:"email"`
	Content string `json:"content"`
	Views int `json:"views"`
	CreatedAt time.Time `json:"time"`
}

var (
	cache *circuitbreaker.StateCache
	profileClient *circuitbreaker.ProtectedHTTPClient
	postsClient   *circuitbreaker.ProtectedHTTPClient
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
		log.Println("[Recommendations] Warning: Redis connection failed: ",err)
	}

	pubsub.Subscribe(cache, serviceName)

	config := circuitbreaker.DefaultConfig()

	profileBreaker := circuitbreaker.NewCircuitBreaker(
		serviceName,
		"profile",
		config,
		cache,
		pubsub,
	)
	profileClient = circuitbreaker.NewProtectedHTTPClient(profileBreaker, config.RequestTimeout)

	postsBreaker := circuitbreaker.NewCircuitBreaker(
		serviceName,
		"posts",
		config,
		cache,
		pubsub,
	)
	postsClient = circuitbreaker.NewProtectedHTTPClient(postsBreaker, config.RequestTimeout)

	log.Printf("[Recommendations] Circuit breakers initialized for profile and posts services")
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
				"state":         state.State.String(),
				"failure_count": state.FailureCount,
				"last_change":   state.LastStateChange,
			}
		}
		ctx.JSON(http.StatusOK, result)
	})

	router.GET("/recommendations", func(ctx *gin.Context) {
		email := ctx.Query("email")
		PROFILESAPI_URL := getEnv("PROFILESAPI_URL", "http://localhost:7000")
		POSTSAPI_URL := getEnv("POSTSAPI_URL", "http://localhost:7001")

		limit := ctx.Query("limit")
		if limit == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "limit not set in query param",
			})
			return
		}

		resp, err := profileClient.Get(PROFILESAPI_URL + "/profiles/similar?email=" + email + "&limit=" + limit)
		if err != nil {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{
				"error": "profile service unavailable: " + err.Error(),
				"circuit": "may be open",
			})
			return
		}
		defer resp.Body.Close()

		var similarProfiles []profile
		err1 := json.NewDecoder(resp.Body).Decode(&similarProfiles)
		if err1 != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err1.Error(),
			})
			return
		}

		recommendedPosts := []post{}
		for _, similarProfile := range similarProfiles {
			resp, err := postsClient.Get(POSTSAPI_URL + "/posts?email=" + similarProfile.Email)
			if err != nil {
				fmt.Printf("[Recommendations] Failed to fetch posts for %s: %v\n", similarProfile.Email, err)
				continue
			}

			var currUserRecommededPosts []post
			jsonerr := json.NewDecoder(resp.Body).Decode(&currUserRecommededPosts)
			resp.Body.Close()
			if jsonerr != nil {
				fmt.Printf("[Recommendations] Failed to parse posts for %s: %v\n",similarProfile.Email, jsonerr)
				continue
			}
			recommendedPosts = append(recommendedPosts, currUserRecommededPosts...)
		}

		ctx.JSON(http.StatusOK, recommendedPosts)
	})

	router.Run(":7003")
}
