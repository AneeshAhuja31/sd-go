package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
	"time"
)

type post struct {
	ID int64 `json:"id"`
	Email string `json:"email"`
	Content string `json:"content"`
	Views int `json:"views"`
	CreatedAt time.Time `json:"time"`
}


func initPQ(hostname string, port int)*sql.DB{
	connStr := "postgres://postgres:pass@" + hostname + ":" + fmt.Sprint(port) + "/postgres?sslmode=disable"
	db,err := sql.Open("postgres",connStr)
	if err != nil {
		fmt.Println("Error connecting to postgres: ",err)
	}
	return db
}

func fetchPostsbyEmail(db*sql.DB,email string)([]post,error){
	query_string := "SELECT id,email,content,views,created_at FROM posts WHERE email = $1"
	rows,err := db.Query(query_string,email)
	if err != nil {
		fmt.Println("Error while querying, ",err)
		return []post{},err
	}
	posts := []post{}
	for rows.Next(){
		var curr_post post
		err := rows.Scan(&curr_post.ID,&curr_post.Email,&curr_post.Content,&curr_post.Views,&curr_post.CreatedAt)
		if err != nil {
			fmt.Println("Error while scanning rows, ",err)
			return []post{},err
		}
		posts = append(posts, curr_post)
	}
	return posts,err
}

func fetchTopPosts(db *sql.DB,limit int)([]post,error){
	query_string := "SELECT id,email,content,views,created_At FROM posts ORDER BY views DESC LIMIT $1"
	rows,err := db.Query(query_string,limit)
	if err != nil {
		fmt.Println("Error in query execution: ",err)
		return []post{},err
	}
	posts := []post{}
	for rows.Next(){
		var curr_post post
		err := rows.Scan(&curr_post.ID,&curr_post.Email,&curr_post.Content,&curr_post.Views,&curr_post.CreatedAt)
		if err != nil {
			fmt.Println("Error while scanning rows, ",err)
			return []post{},err
		}
		posts = append(posts, curr_post)
	}
	return posts,err
}