package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type profile struct {
	ID int64 `json:"id"`
	Email string `json:"email"`
	Username string `json:"username"`
	DOB time.Time `json:"dob"`
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

func main(){
	godotenv.Load()
	router := gin.Default()
	router.GET("/recommendations",func(ctx *gin.Context) {
		email := ctx.Query("email")
		PROFILESAPI_URL,ok := os.LookupEnv("PROFILESAPI_URL")
		if !ok {
			PROFILESAPI_URL = "http://localhost:7000"
		}
		POSTSAPI_URL,ok := os.LookupEnv("POSTSAPI_URL")
		if !ok {
			POSTSAPI_URL = "http://localhost:7001"
		}
		// resp,err := http.Get(PROFILESAPI_URL+"/profile?email="+email)
		// if err != nil {
		// 	ctx.JSON(http.StatusInternalServerError,gin.H{
		// 		"error":err,
		// 	})
		// 	return
		// }
		// var user_profile profile
		// err1 := json.NewDecoder(resp.Body).Decode(&user_profile)
		// if err1 != nil {
		// 	ctx.JSON(http.StatusInternalServerError,gin.H{
		// 		"error":err1,
		// 	})
		// 	return
		// }
		limit := ctx.Query("limit")

		resp,err := http.Get(PROFILESAPI_URL + "/profiles/similar?email="+email+"&limit="+limit)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{
				"error":err.Error(),
			})
			return
		}
		defer resp.Body.Close()
		var similarProfiles []profile
		err1 := json.NewDecoder(resp.Body).Decode(&similarProfiles)
		if err1 != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{
				"error":err1,
			})
			return
		}
		recommendedPosts := []post{}
		for _,similarProfile := range(similarProfiles){
			resp,err := http.Get(POSTSAPI_URL+"/posts?email="+similarProfile.Email)
			if err != nil{
				ctx.JSON(http.StatusInternalServerError,gin.H{
					"error":err.Error(),
				})
				return
			}
			
			var currUserRecommededPosts []post
			jsonerr := json.NewDecoder(resp.Body).Decode(&currUserRecommededPosts)
			resp.Body.Close()
			if jsonerr != nil {
				ctx.JSON(http.StatusInternalServerError,gin.H{
					"error":jsonerr.Error(),
				})
				return
			}
			recommendedPosts = append(recommendedPosts, currUserRecommededPosts...)
		}
		ctx.JSON(http.StatusOK, recommendedPosts)
	})
	router.Run(":7003")
}