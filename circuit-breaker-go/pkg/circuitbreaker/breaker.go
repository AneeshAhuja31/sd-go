package circuitbreaker

import (
	"errors"
	"log"
	"sync"
	"time"
)

var ErrCircuitOpen = errors.New("circuit breaker is open")
var ErrTooManyRequests = errors.New("too many requests in half-open state")
type CircuitBreaker struct {
	mu sync.Mutex
	config Config
	cache *StateCache
	pubsub *RedisPubSub
	serviceName string
	targetService string

	halfOpenRequests int
}

func NewCircuitBreaker(serviceName string,targetService string,config Config,cache *StateCache,pubsub *RedisPubSub) *CircuitBreaker {
	return &CircuitBreaker{
		config: config,
		cache:cache,
		pubsub: pubsub,
		serviceName: serviceName,
		targetService: targetService,
	}
}

func (cb *CircuitBreaker) CanExecute() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	state := cb.cache.GetorCreate(cb.serviceName,cb.targetService)
	switch state.State{
		case StateClosed:
			return nil
		case StateOpen:
			if time.Now().After(state.OpenUntil) { //check if timeot is over
				cb.transitionTo(StateHalfOpen,"time elapsed")
				cb.halfOpenRequests = 0
				return nil
			}
			return ErrCircuitOpen
		case StateHalfOpen:
			if cb.halfOpenRequests >= cb.config.HalfOpenMaxRequests {
				return ErrTooManyRequests
			}
			cb.halfOpenRequests++
			return nil
		}
	return nil
}

func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	state := cb.cache.GetorCreate(cb.serviceName,cb.targetService)
	switch state.State {
		case StateClosed:
			state.FailureCount = 0 
		case StateHalfOpen:
			state.SuccessCount ++
			if state.SuccessCount >= cb.config.SuccessThreshold{
				cb.transitionTo(StateClosed,"success threshold reached")
			}
	}
}

func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	state := cb.cache.GetorCreate(cb.serviceName, cb.targetService)
	state.FailureCount ++
	state.LastFailureTime = time.Now()

	switch state.State{
	case StateClosed:
		if state.FailureCount >= cb.config.FailureThreshold{
			state.OpenUntil = time.Now().Add(cb.config.Timeout)
			cb.transitionTo(StateOpen,"failure threshold reached")
		}
	case StateHalfOpen:
		state.OpenUntil = time.Now().Add(cb.config.Timeout)
		cb.transitionTo(StateOpen,"failure in half-open")
	}
}

func (cb *CircuitBreaker) transitionTo(newState State,reason string){
	state := cb.cache.GetorCreate(cb.serviceName,cb.targetService)
	oldState := state.State

	state.State = newState
	state.LastStateChange = time.Now()

	if newState == StateClosed {
		state.FailureCount = 0
		state.SuccessCount = 0
	} else if newState == StateHalfOpen {
		state.SuccessCount = 0
	}

	log.Printf("[CircuitBreaker] %s:%s transitioned %s -> %s (%s)",cb.serviceName, cb.targetService, oldState.String(), newState.String(), reason)
	msg := StateChangeMessage{
		SourceService: cb.serviceName,
		TargetService: cb.targetService,
		OldState:      oldState,
		NewState:      newState,
		FailureCount:  state.FailureCount,
		Timestamp:     time.Now(),
		Reason:        reason,
	}

	if err := cb.pubsub.Publish(msg); err != nil {
		log.Printf("[CircuitBreaker] Failed to publish state change: %v", err)
	}
}

func (cb *CircuitBreaker) GetState() *CircuitBreakerState {
	return cb.cache.GetorCreate(cb.serviceName,cb.targetService)
}