package main

import (
	"fmt"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)


func main(){
	db := initPQ("localhost",5432)
	router := gin.Default()
	router.GET("/profile", func(c*gin.Context){
		email := c.Query("email")
		if email == "" {
			c.JSON(400, gin.H{
				"error":"Empty email",
			})
			return
		}
		userProfile,err := fetchProfileData(db,email)
		if err != nil {
			fmt.Println("Error: ",err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK,userProfile)
	})

	router.GET("/profiles/similar",func(ctx *gin.Context) {
		email := ctx.Query("email")
		if email == "" {
			ctx.JSON(400, gin.H{
				"error":"Empty email",
			})
			return
		}
		userProfile,err := fetchProfileData(db,email)
		if err != nil {
			fmt.Println("Error: ",err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":err.Error(),
			})
			return
		}
		limit := ctx.Query("limit")
		num_limit,err := strconv.Atoi(limit)
		if err != nil {
			ctx.JSON(400, gin.H{
				"error":err.Error(),
			})
			return
		}
		similarUsers, err := fetchProfilesByHobbies(db,userProfile.Hobbies,email, num_limit)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{
				"error":err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK,similarUsers)
	})

	err := router.Run(":7000")
	handleError(err)
}
