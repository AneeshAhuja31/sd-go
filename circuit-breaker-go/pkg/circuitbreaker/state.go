package circuitbreaker

import "time"

type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

func (s State) String() string {
	switch s {
	case StateClosed:
		return "CLOSED"
	case StateOpen:
		return "OPEN"
	case StateHalfOpen:
		return "HALF_OPEN"
	default:
		return "UNKNOWN"
	}
}

type CircuitBreakerState struct {
	ServiceName string `json:"service_name"`
	TargetService string `json:"target_service"`
	State State `json:"state"`
	FailureCount int `json:"failure_count"`
	SuccessCount int `json:"success_count"`
	LastFailureTime time.Time `json:"last_failure_time"`
	LastStateChange time.Time `json:"last_state_change"`
	OpenUntil time.Time `json:"open_until"`
}

