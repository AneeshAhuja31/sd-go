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



// func makeNode(slot int)*Node{
// 	return &Node{
// 		ID: "node-"+fmt.Sprint(port),
// 	}
// }