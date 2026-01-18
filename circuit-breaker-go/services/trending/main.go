package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type post struct {
	ID int64 `json64:"id"`
	Email string `json:"email"`
	Content string `json:"content"`
	Views int `json:"views"`
	CreatedAt time.Time `json:"created_at"`
}

func main(){
	godotenv.Load()
	router := gin.Default()
	router.GET("/trending",func(ctx *gin.Context) {
		POSTSAPI_URL,ok := os.LookupEnv("POSTSAPI_URL")
		if !ok {
			POSTSAPI_URL = "http://localhost:7001"
		}
		resp,err := http.Get(POSTSAPI_URL+"/topPosts?limit=5")
		if err != nil{
			ctx.JSON(http.StatusServiceUnavailable,gin.H{
				"error":"invalid response from posts service",
			})
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error parsing request body: ",err)
			ctx.JSON(http.StatusInternalServerError,gin.H{
				"error":err,
			})
			return
		}
		var posts []post
		json.Unmarshal(body, &posts)
		
		ctx.JSON(http.StatusOK,gin.H{
			"trending":posts,
		})
	})
	router.Run(":7002")
}