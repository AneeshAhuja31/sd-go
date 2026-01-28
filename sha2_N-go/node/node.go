package node

import (
	"net/http"
	"database/sql"
	"fmt"
	"sha-go/hash"
)

type Node struct {
	ID string
	Slot int
	Hash uint64
	HttpServer *http.Server
	DB *sql.DB
}



func MakeNode(port int,httpaddr *http.Server,db *sql.DB,n int)*Node{
	id := "node-"+fmt.Sprint(port)
	hashVal := hash.Hash(id)
	slot := hash.GetSlot(hashVal,n)
	return &Node{
		ID: id,
		Slot: slot,
		Hash: hashVal,
		HttpServer: httpaddr,
		DB: db,
	}
}