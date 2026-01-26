package node

import (
	"net/http"
	"database/sql"
	"fmt"
)

type Node struct {
	ID string
	slot int
	hash uint64
	httpAddr *http.Server
	db *sql.DB
}



func makeNode(port int,slot int,hash uint64,httpaddr *http.Server,db *sql.DB)*Node{
	return &Node{
		ID: "node-"+fmt.Sprint(port),
		slot: slot,
		hash: hash,
		httpAddr: httpaddr,
		db: db,
	}
}