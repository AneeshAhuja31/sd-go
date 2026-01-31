package main

import (
	"context"
	"fmt"
	"leaky-bucket-rate-limiter-go/config"
	"log"
	"net/http"

)


func RateLimitMiddleware(bucket *LeakyBucket, ctx context.Context, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		key := fmt.Sprintf("rate_limit:%s", ip)

		allowed, err := bucket.Allow(ctx, key)
		if err != nil {
			log.Printf("rate limiter error: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if !allowed {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	rdb, ctx := config.InitRedisClient("localhost", 6379)
	defer rdb.Close()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	log.Println("connected to redis")

	//capacity=10 requests, leakRate=2 requests/sec
	bucket := NewLeakyBucket(rdb, 10, 2)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "request successful")
	})

	handler := RateLimitMiddleware(bucket, ctx, mux)

	log.Println("server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

