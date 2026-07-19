package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

type State int

const (
	Closed State = iota
	Open
	HalfOpen
)

type CircuitBreaker struct {
	state        State
	failureCount int
	maxFailures  int
	lastFailure  time.Time
	timeout      time.Duration
	mu           sync.Mutex
}

func New(maxFailures int, timeout time.Duration) *CircuitBreaker {

	newCircuitBreaker := &CircuitBreaker{maxFailures: maxFailures, timeout: timeout}
	return newCircuitBreaker

}

func (c *CircuitBreaker) Execute(fn func() error) error {

	c.mu.Lock()
	if c.state == Open && time.Since(c.lastFailure) < c.timeout {
		c.mu.Unlock()
		return errors.New("Problem occured processing the task")
	}
	if c.state == Open && time.Since(c.lastFailure) >= c.timeout {
		c.state = HalfOpen
	}
	c.mu.Unlock()
	err := fn()

	c.mu.Lock()
	if err != nil {

		c.failureCount++

		if c.failureCount >= c.maxFailures {
			c.state = Open
			c.lastFailure = time.Now()
		}
	}
	if err == nil {

		c.failureCount = 0
		c.state = Closed
	}
	c.mu.Unlock()
	return err
}
