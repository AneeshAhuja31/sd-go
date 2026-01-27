package node

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
)

func ShutdownNode(node *Node, dbHost string, dbPort int) error {
	log.Printf("Shutting down node %s", node.ID)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := node.HttpServer.Shutdown(ctx); err != nil {
		log.Printf("Error shutting down HTTP server for node %s: %v", node.ID, err)
	}

	dbName := strings.Replace(node.ID, "-", "_", -1)

	if err := node.DB.Close(); err != nil {
		log.Printf("Error closing database for node %s: %v", node.ID, err)
		return err
	}

	connStr := fmt.Sprintf("postgres://postgres:pass@%s:%d/postgres?sslmode=disable", dbHost, dbPort)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Error connecting to postgres for cleanup: %v", err)
		return err
	}
	defer db.Close()

	dropQuery := fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName)
	_, err = db.Exec(dropQuery)
	if err != nil {
		log.Printf("Error dropping database %s: %v", dbName, err)
		return err
	}

	log.Printf("Node %s shut down successfully and database %s deleted", node.ID, dbName)
	return nil
}
