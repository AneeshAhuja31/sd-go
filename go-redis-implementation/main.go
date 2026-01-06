package main

import (
	"log"
	"net/http"
	"github.com/redis/go-redis/v9"
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"encoding/json"
)

type jsonBodyStruct struct {
	UserList json.RawMessage `json:"userlist"`
	Location string `json:"location"`
}

func fetchData(db *sql.DB,rdb *redis.Client,ctx context.Context)[]byte{
	var users []byte
	users,err := fetchDataRedis(rdb,ctx)
	var location string
	if err != nil {
		log.Println("GOT USERS FROM POSTGRES")
		location = "postgres"
		users = fetchDataPq(db)
		setDataRedis(users,rdb,ctx)
	} else {
		log.Println("GOT USERS FROM REDIS")
		location = "redis"
	}
	jsonbodystruct := jsonBodyStruct{
		UserList: json.RawMessage(users),
		Location: location,
	}
	jsonbody,err := json.Marshal(jsonbodystruct)
	if err != nil {
		panic(err)
	}
	return jsonbody
}

func main(){
	db := connectPq("localhost",5432)
	rdb,ctx := connectRedis("localhost",6379)
	http.HandleFunc("/fetch-users",func(w http.ResponseWriter, r *http.Request) {
		users := fetchData(db,rdb,ctx)
		w.Header().Set("Content-Type","application/json")
		w.Write(users)
	})

	log.Println("Server running on 8080")
	http.ListenAndServe(":8080",nil)
}