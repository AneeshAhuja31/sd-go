package api

import (
	"database/sql"
	"net/http"
	"github.com/gin-gonic/gin"
	"url-shortener-go/internal/cache"
	"url-shortener-go/internal/ticket"
)

type API struct {
	DB *sql.DB
	Cache *cache.RedisCache
	Ticket *ticket.LocalTicketClient
	Events chan string
}

func NewAPI(db *sql.DB, redis *cache.RedisCache, ticketClient *ticket.LocalTicketClient, events chan string) *API {
	return &API{
		DB: db,
		Cache: redis,
		Ticket: ticketClient,
		Events: events,
	}
}

func (a *API) RegisterRoutes(r *gin.Engine) {

	r.POST("/shorten", func(c *gin.Context) {
		var req struct {
			LongURL string `json:"long_url"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid request"})
			return
		}
		id, _ := a.Ticket.NextID()
		code := ticket.EncodeBase62(id)

		_, err := a.DB.Exec(
			"INSERT INTO urls(short_code, long_url) VALUES($1,$2)",
			code, req.LongURL,
		)

		if err != nil {
			c.JSON(500, gin.H{"error": "db insert failed"})
			return
		}

		c.JSON(200, gin.H{
			"short_url": "http://localhost:8080/" + code,
		})
	})
	//redirect
	r.GET("/:code", func(c *gin.Context) {

		code := c.Param("code")

		longURL, err := a.Cache.Get(code)
		if err == nil {
			a.Events <- code
			c.Redirect(http.StatusFound, longURL)
			return
		}
		err = a.DB.QueryRow(
			"SELECT long_url FROM urls WHERE short_code=$1",
			code,
		).Scan(&longURL)

		if err != nil {
			c.JSON(404, gin.H{"error": "not found"})
			return
		}

		a.Cache.Set(code, longURL)

		a.Events <- code

		c.Redirect(http.StatusFound, longURL)
	})
}
