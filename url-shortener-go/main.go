package main

import (
	"context"
	"log"
	"github.com/gin-gonic/gin"
	"url-shortener/internal/api"
	"url-shortener/internal/cache"
	"url-shortener/internal/db"
	"url-shortener/internal/ticket"
	"url-shortener/internal/worker"
)

func main() {

	clickEvents := make(chan string, 100)

	pg := db.InitDB("localhost",5432)

	ctx := context.Background()
	redis, err := cache.InitRedis(ctx, "localhost", 6379)
	if err != nil {
		log.Fatal(err)
	}

	ts := ticket.NewTicketServer(1, 1000, 100)
	client := ticket.NewLocalTicketClient(ts)

	worker.StartConsumer(pg, clickEvents)

	r := gin.Default()

	apiServer := api.NewAPI(pg, redis, client, clickEvents)
	apiServer.RegisterRoutes(r)

	r.Run(":8080")
}
