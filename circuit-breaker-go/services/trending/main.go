package main

import (
	"encoding/json"
	"fmt"

	//"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type post struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Content   string    `json:"content"`
	Views     int       `json:"views"`
	CreatedAt time.Time `json:"created_at"`
}

func main() {
	godotenv.Load()
	router := gin.Default()
	router.GET("/trending", func(ctx *gin.Context) {
		POSTSAPI_URL, ok := os.LookupEnv("POSTSAPI_URL")
		if !ok {
			POSTSAPI_URL = "http://localhost:7001"
		}
		limit := ctx.Query("limit")
		if limit == ""{
			ctx.JSON(http.StatusBadRequest,gin.H{
				"error":"limit not set in query param",
			})
		}
		resp, err := http.Get(POSTSAPI_URL + "/topPosts?limit=" + limit)
		if err != nil {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{
				"error": "invalid response from posts service",
			})
			return
		}
		// defer resp.Body.Close()

		var posts []post
		err1 := json.NewDecoder(resp.Body).Decode(&posts)
		if err1 != nil {
			fmt.Println("Error parsing json: ", err1)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err1.Error(),
			})
			return
		}

		// body, err := io.ReadAll(resp.Body)
		// if err != nil {
		// 	fmt.Println("Error parsing request body: ",err)
		// 	ctx.JSON(http.StatusInternalServerError,gin.H{
		// 		"error":err,
		// 	})
		// 	return
		// }
		// // var posts []post
		// json.Unmarshal(body, &posts)

		ctx.JSON(http.StatusOK, posts)
	})
	router.Run(":7002")
}
