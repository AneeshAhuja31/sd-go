package main

import (
	_ "github.com/lib/pq"
	"fmt"
	"database/sql"
	"log"
	"encoding/json"
)

type User struct {
	ID int `json:"id"`
	UserName string `json:"username"`
}

func connectPq(host string, port int) *sql.DB{
	connStr := "postgres://postgres:pass@" + host + ":" + fmt.Sprint(port) + "/postgres?sslmode=disable"
	db,err := sql.Open("postgres",connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func fetchDataPq(db *sql.DB) []byte{
	rows,err := db.Query("SELECT * FROM users")
	if err != nil {
		panic(err.Error())
	}
	var users [] User

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID,&user.UserName)
		if err != nil {
			panic(err.Error())
		}
		users = append(users, user)
	}
	data,err := json.Marshal(users)
	if err != nil {
		panic(err)
	}

	return data
}