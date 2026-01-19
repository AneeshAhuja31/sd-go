package circuitbreaker

import (
	"fmt"
	"net/http"
	"time"
)

type ProtectedHTTPClient struct {
	client  *http.Client
	breaker *CircuitBreaker
}

func NewProtectedHTTPClient(breaker *CircuitBreaker,timeout time.Duration) *ProtectedHTTPClient {
	return &ProtectedHTTPClient{
		client: &http.Client{
			Timeout: timeout,
		},
		breaker: breaker,
	}
}

func (p *ProtectedHTTPClient) Get(url string) (*http.Response, error) {
	//check if circuit allows execution
	if err := p.breaker.CanExecute(); err != nil {
		return nil, fmt.Errorf("circuit breaker: %w", err)
	}
	resp, err := p.client.Get(url)

	if err != nil {
		p.breaker.RecordFailure()
		return nil, err
	}

	if resp.StatusCode >= 500 {
		p.breaker.RecordFailure()
		return resp, fmt.Errorf("server error: status %d", resp.StatusCode)
	}

	p.breaker.RecordSuccess()
	return resp, nil
}

//to perform other request types w payload not needed rn
func (p *ProtectedHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if err := p.breaker.CanExecute(); err != nil {
		return nil, fmt.Errorf("circuit breaker: %w", err)
	}

	resp, err := p.client.Do(req)

	if err != nil {
		p.breaker.RecordFailure()
		return nil, err
	}

	if resp.StatusCode >= 500 {
		p.breaker.RecordFailure()
		return resp, fmt.Errorf("server error: status %d", resp.StatusCode)
	}

	p.breaker.RecordSuccess()
	return resp, nil
}
