package circuitbreaker

import (
	"sync"
	"time"
)

type CacheKey string

func NewCacheKey(sourceService, targetService string)CacheKey{
	return CacheKey(sourceService+":"+targetService)
}

type StateCache struct {
	mu sync.RWMutex
	states map[CacheKey]*CircuitBreakerState
}

func NewStateCache() *StateCache {
	return &StateCache{
		states: make(map[CacheKey]*CircuitBreakerState),
	}
}

func (c *StateCache) Set(key CacheKey, state *CircuitBreakerState) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.states[key] = state
}

func (c *StateCache) GetorCreate(sourceService string, targetService string) *CircuitBreakerState{
	key := NewCacheKey(sourceService,targetService)
	c.mu.Lock()
	defer c.mu.Unlock()
	state,ok := c.states[key]
	if ok {
		return state
	}
	state = &CircuitBreakerState{
		ServiceName: sourceService,
		TargetService: targetService,
		State: StateClosed,
		FailureCount: 0,
		SuccessCount: 0,
		LastStateChange: time.Now(),
	}
	c.states[key] = state
	return state
}

func (c *StateCache) GetAll() map[CacheKey]*CircuitBreakerState {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make(map[CacheKey]*CircuitBreakerState)
	for k, v := range c.states {
		result[k] = v
	}
	return result
}

