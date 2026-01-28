package coord

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sha-go/hash"
	"sha-go/node"
	"sha-go/ring"
	"github.com/redis/go-redis/v9"
)

const (
	NodesKey = "sha2n:nodes"
)

type NodeInfo struct {
	ID string `json:"id"`
	Slot int `json:"slot"`
	Hash uint64 `json:"hash"`
	HttpAddr string `json:"http_addr"`
}

func RegisterNode(shaRing *ring.Ring, n *node.Node) error {
	nodeInfo := NodeInfo{
		ID: n.ID,
		Slot: n.Slot,
		Hash: n.Hash,
		HttpAddr: n.HttpServer.Addr,
	}

	nodeJSON, err := json.Marshal(nodeInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal node info: %w", err)
	}

	err = shaRing.RedisConn.HSet(shaRing.Ctx, NodesKey, n.ID, nodeJSON).Err()
	if err != nil {
		return fmt.Errorf("failed to register node in Redis: %w", err)
	}

	shaRing.Ring[n.Slot] = *n

	log.Printf("Registered node %s at slot %d", n.ID, n.Slot)
	return nil
}

func DeregisterNode(shaRing *ring.Ring, nodeID string, slot int) error {
	err := shaRing.RedisConn.HDel(shaRing.Ctx, NodesKey, nodeID).Err()
	if err != nil {
		return fmt.Errorf("failed to deregister node from Redis: %w", err)
	}

	shaRing.Ring[slot] = node.Node{}

	log.Printf("Deregistered node %s from slot %d", nodeID, slot)
	return nil
}

func GetAllNodes(ctx context.Context, redisClient *redis.Client) ([]NodeInfo, error) {
	nodesMap, err := redisClient.HGetAll(ctx, NodesKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes from Redis: %w", err)
	}

	var nodes []NodeInfo
	for _, nodeJSON := range nodesMap {
		var nodeInfo NodeInfo
		err := json.Unmarshal([]byte(nodeJSON), &nodeInfo)
		if err != nil {
			log.Printf("Failed to unmarshal node info: %v", err)
			continue
		}
		nodes = append(nodes, nodeInfo)
	}

	return nodes, nil
}

func LoadRingFromRedis(shaRing *ring.Ring) error {
	nodes, err := GetAllNodes(shaRing.Ctx, shaRing.RedisConn)
	if err != nil {
		return err
	}

	shaRing.Ring = make([]node.Node,hash.TotalSlots)

	for _,nodeInfo := range nodes {
		n := node.Node{
			ID: nodeInfo.ID,
			Slot: nodeInfo.Slot,
			Hash: nodeInfo.Hash,
		}
		shaRing.Ring[nodeInfo.Slot] = n
	}

	log.Printf("Loaded %d nodes from Redis into ring with %d slots", len(nodes), hash.TotalSlots)
	return nil
}
