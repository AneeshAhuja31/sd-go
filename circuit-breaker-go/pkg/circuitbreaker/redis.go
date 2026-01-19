package circuitbreaker

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
	"fmt"
)

type StateChangeMessage struct {
	SourceService string `json:"source_service"`
	TargetService string `json:"target_service"`
	OldState State `json:"old_state"`
	NewState State `json:"new_state"`
	FailureCount int `json:"failure_count"`
	Timestamp time.Time `json:"timestamp"`
	Reason string `json:"reason"`
}

type RedisPubSub struct {
	client *redis.Client
	channel string
	ctx context.Context
}

func NewRedisPubSub(config RedisConfig) *RedisPubSub {
	client := redis.NewClient(&redis.Options{
		Addr: config.Addr,
		Password: config.Password,
		DB:config.DB,
	})

	return &RedisPubSub{
		client: client,
		channel: config.Channel,
		ctx: context.Background(),
	}
}


func (r *RedisPubSub) Ping() error {
	return r.client.Ping(r.ctx).Err()
}

func (r *RedisPubSub) Publish(msg StateChangeMessage) error {
	data,err := json.Marshal(msg)
	if err != nil{
		return fmt.Errorf("failed to marshall message")
	}
	err = r.client.Publish(r.ctx,r.channel,string(data)).Err()
	if err != nil {
		return fmt.Errorf("failed to publish message")
	}

	log.Printf("[REDIS] Published state change: %s:%s %s -> %s",msg.SourceService,msg.TargetService,msg.OldState.String(),msg.NewState.String())
	return nil
}

func (r* RedisPubSub) Subscribe(cache *StateCache, serviceName string){
	sub := r.client.Subscribe(r.ctx,r.channel)
	ch := sub.Channel()

	log.Printf("[Redis] %s subscribed to channel: %s", serviceName, r.channel)
	go func() {
		for msg := range ch {
			var stateMsg StateChangeMessage
			if err := json.Unmarshal([]byte(msg.Payload), &stateMsg); err != nil {
				log.Printf("[Redis] Failed to unmarshal message: %v", err)
				continue
			}

			key := NewCacheKey(stateMsg.SourceService, stateMsg.TargetService) //update local cache with recevied state
			state := &CircuitBreakerState{
				ServiceName:     stateMsg.SourceService,
				TargetService:   stateMsg.TargetService,
				State:           stateMsg.NewState,
				FailureCount:    stateMsg.FailureCount,
				LastStateChange: stateMsg.Timestamp,
			}

			cache.Set(key, state)
			log.Printf("[Redis] %s updated cache: %s -> %s",
				serviceName, key, stateMsg.NewState.String())
		}
	}()
}

func (r *RedisPubSub) Close() error {
	return r.client.Close()
}