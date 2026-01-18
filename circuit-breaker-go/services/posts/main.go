package main

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"net/http"
	"strconv"
)

func main(){
	db := initPQ("localhost",5432)
	router := gin.Default()
	router.GET("/posts",func(ctx *gin.Context) {
		email := ctx.Query("email")
		if email == "" {
			ctx.JSON(400,gin.H{
				"error":"Email empty",
			})
			return
		}
		posts_data,err := fetchPostsbyEmail(db,email)
		if err != nil{
			ctx.JSON(http.StatusInternalServerError,gin.H{
				"error":err.Error(),
			})
		}
		ctx.JSON(http.StatusOK,posts_data)
	})
	
	router.GET("/topPosts",func(ctx *gin.Context) {
		limit := ctx.Query("limit")
		num_limit,err := strconv.Atoi(limit)
		if err != nil{
			fmt.Println("Invalid (non interger) character in limit url query: ",err)
			ctx.JSON(http.StatusBadRequest,gin.H{
				"error":err.Error(),
			})
			return
		}
		trendingposts, err := fetchTopPosts(db,num_limit)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{
				"error":err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK,trendingposts)

	})
	err := router.Run(":7001")
	if err != nil{
		fmt.Println("Error running server: ", err)
	}
}