// Package semaphore provides a weighted concurrency limiter that restricts
// the number of goroutines executing a critical section simultaneously.
package semaphore

import (
	"context"
	"fmt"
)

// Semaphore controls access to a resource pool with a fixed number of slots.
type Semaphore struct {
	ch chan struct{}
}

// New creates a Semaphore with the given capacity.
// It returns an error if cap is less than 1.
func New(cap int) (*Semaphore, error) {
	if cap < 1 {
		return nil, fmt.Errorf("semaphore: capacity must be at least 1, got %d", cap)
	}
	return &Semaphore{ch: make(chan struct{}, cap)}, nil
}

// Acquire blocks until a slot is available or ctx is cancelled.
// Returns ctx.Err() if the context is done before a slot is acquired.
func (s *Semaphore) Acquire(ctx context.Context) error {
	select {
	case s.ch <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// TryAcquire attempts to acquire a slot without blocking.
// Returns true if a slot was acquired, false otherwise.
func (s *Semaphore) TryAcquire() bool {
	select {
	case s.ch <- struct{}{}:
		return true
	default:
		return false
	}
}

// Release returns a slot to the semaphore.
// It panics if called more times than Acquire.
func (s *Semaphore) Release() {
	select {
	case <-s.ch:
	default:
		panic("semaphore: Release called without a matching Acquire")
	}
}

// Available returns the number of slots currently free.
func (s *Semaphore) Available() int {
	return cap(s.ch) - len(s.ch)
}

// Cap returns the total capacity of the semaphore.
func (s *Semaphore) Cap() int {
	return cap(s.ch)
}
