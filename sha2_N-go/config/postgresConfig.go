package config 

import (
	_ "github.com/lib/pq"
	"database/sql"
	"fmt"
	"log"
)

type FileRecord struct{
	Key string
	Value string
	Hash uint64
}

func InitPQ(host string, port int)*sql.DB{
	connStr := "postgres://postgres:pass@" + host + ":" + fmt.Sprint(port) + "/postgres?sslmode=disable"
	db, err := sql.Open("postgres",connStr)
	if err!=nil{
		log.Fatal(err)
	}
	log.Println("Setup pq connection")
	return db
}

