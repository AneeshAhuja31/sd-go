package node

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sha-go/config"
)

func InitNode(port int, dbHost string, dbPort int,n int) (*Node, error) {
	db := config.InitPQ(dbHost, dbPort)
	dbName := fmt.Sprintf("node_%d", port)
	dbCreationStatement := fmt.Sprintf("CREATE DATABASE %s",dbName)
	_,err := db.Exec(dbCreationStatement)
	if err != nil {
		log.Printf("Database %s already exist: %v", dbName, err)
	}
	db.Close()

	connStr := fmt.Sprintf("postgres://postgres:pass@%s:%d/%s?sslmode=disable", dbHost, dbPort, dbName)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database %s: %w", dbName, err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	httpServer := &http.Server{Addr: fmt.Sprintf("localhost:%d", port)}

	node := MakeNode(port, httpServer, db, n)

	createTableQuery := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS files_%d (
			key TEXT PRIMARY KEY,
			value TEXT,
			hash BIGINT
		)
	`, node.Slot)

	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	log.Printf("Initialized node %s at slot %d on %s", node.ID, node.Slot, httpServer.Addr)
	return node, nil
}

func StartNode(node *Node, handler http.Handler) error {
	node.HttpServer.Handler = handler

	log.Printf("Starting node %s HTTP server on %s", node.ID, node.HttpServer.Addr)

	go func() {
		if err := node.HttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error for node %s: %v", node.ID, err)
		}
	}()

	return nil
}
