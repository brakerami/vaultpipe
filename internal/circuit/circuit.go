// Package circuit implements a simple circuit breaker for protecting
// downstream calls (e.g. Vault) from cascading failures.
package circuit

import (
	"errors"
	"sync"
	"time"
)

// State represents the current circuit breaker state.
type State int

const (
	StateClosed   State = iota // normal operation
	StateOpen                  // blocking calls
	StateHalfOpen              // testing recovery
)

// ErrOpen is returned when the circuit is open and calls are rejected.
var ErrOpen = errors.New("circuit breaker is open")

// Breaker is a simple three-state circuit breaker.
type Breaker struct {
	mu           sync.Mutex
	state        State
	failures     int
	threshold    int
	resetTimeout time.Duration
	openedAt     time.Time
}

// New creates a Breaker that opens after threshold consecutive failures
// and attempts recovery after resetTimeout.
func New(threshold int, resetTimeout time.Duration) (*Breaker, error) {
	if threshold <= 0 {
		return nil, errors.New("circuit: threshold must be greater than zero")
	}
	if resetTimeout <= 0 {
		return nil, errors.New("circuit: resetTimeout must be greater than zero")
	}
	return &Breaker{
		threshold:    threshold,
		resetTimeout: resetTimeout,
	}, nil
}

// Allow reports whether the call should proceed. It transitions
// an open circuit to half-open once the reset timeout has elapsed.
func (b *Breaker) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	switch b.state {
	case StateClosed:
		return true
	case StateOpen:
		if time.Since(b.openedAt) >= b.resetTimeout {
			b.state = StateHalfOpen
			return true
		}
		return false
	case StateHalfOpen:
		return true
	}
	return false
}

// RecordSuccess resets the breaker to closed on a successful call.
func (b *Breaker) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
	b.state = StateClosed
}

// RecordFailure increments the failure count and opens the circuit
// if the threshold is reached.
func (b *Breaker) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures++
	if b.failures >= b.threshold {
		b.state = StateOpen
		b.openedAt = time.Now()
	}
}

// State returns the current breaker state.
func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}
