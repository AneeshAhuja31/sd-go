package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"sha-go/api"
	"sha-go/coord"
	"sha-go/node"
	"sha-go/ring"
)

const (
	N = 16 
	DefaultDBHost = "localhost"
	DefaultDBPort = 5432
)

func main() {
	port := flag.Int("port", 8080, "HTTP server port for this node")
	flag.Parse()

	log.Println("Starting SHA-2^N distributed system node on port ", *port)

	shaRing := ring.MakeRing()
	err := coord.LoadRingFromRedis(shaRing)
	if err != nil {
		log.Println("Warning: Could not load ring from Redis: ", err)
	}

	n, err := node.InitNode(*port, DefaultDBHost, DefaultDBPort, N)
	if err != nil {
		log.Fatalf("Failed to initialize node: %v", err)
	}

	err = coord.RegisterNode(shaRing, n)
	if err != nil {
		log.Fatalf("Failed to register node: %v", err)
	}

	srv := &api.Server{
		Node: n,
		Ring: shaRing,
	}
	router := api.SetupRouter(srv)

	err = node.StartNode(n, router)
	if err != nil {
		log.Fatalf("Failed to start node: %v",err)
	}

	log.Printf("Node %s is running at slot %d on %s",n.ID,n.Slot,n.HttpServer.Addr)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutdown signal received, cleaning up!")

	err = coord.DeregisterNode(shaRing, n.ID, n.Slot)
	if err != nil {
		log.Println("Error deregistering node: ", err)
	}
 
	err = node.ShutdownNode(n, DefaultDBHost, DefaultDBPort)
	if err != nil {
		log.Println("Error during shutdown:", err)
	}
	log.Println("Node shut down successfully")
}
