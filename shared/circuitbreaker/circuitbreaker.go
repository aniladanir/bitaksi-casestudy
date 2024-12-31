package circuitbreaker

import (
	"fmt"
	"sync"
	"time"
)

type CircuitBreakerState int

const (
	StateClosed CircuitBreakerState = iota
	StateOpen
	StateHalfOpen
)

type CircuitBreaker struct {
	state        CircuitBreakerState
	failureCount int
	lastAttempt  time.Time
	mu           sync.RWMutex
	maxFailures  int
	retryTimeout time.Duration
}

type CircuitBreakerOption func(*CircuitBreaker)

func WithMaxFailures(maxFailures int) CircuitBreakerOption {
	return func(cb *CircuitBreaker) {
		cb.maxFailures = maxFailures
	}
}
func WithRetryTimeout(retryTimeout time.Duration) CircuitBreakerOption {
	return func(cb *CircuitBreaker) {
		cb.retryTimeout = retryTimeout
	}
}

func NewCircuitBreaker(opts ...CircuitBreakerOption) *CircuitBreaker {
	cb := &CircuitBreaker{
		state:        StateClosed,
		maxFailures:  3,
		retryTimeout: 15 * time.Second,
	}
	for _, opt := range opts {
		opt(cb)
	}
	return cb
}

func (cb *CircuitBreaker) Execute(requestFunc func() (interface{}, error)) (interface{}, error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.state == StateOpen {
		if time.Since(cb.lastAttempt) < cb.retryTimeout {
			return nil, fmt.Errorf("circuit breaker is open")
		}
		cb.state = StateHalfOpen
	}

	result, err := requestFunc()
	if err != nil {
		cb.handleFailure()
		return nil, err
	}
	cb.handleSuccess()
	return result, nil
}

func (cb *CircuitBreaker) handleSuccess() {
	cb.failureCount = 0
	if cb.state == StateHalfOpen {
		cb.state = StateClosed
		fmt.Println("Circuit breaker closed")
	}
}

func (cb *CircuitBreaker) handleFailure() {
	cb.failureCount++
	if cb.failureCount >= cb.maxFailures {
		cb.state = StateOpen
		cb.lastAttempt = time.Now()
		fmt.Println("Circuit breaker opened")
	}
}
