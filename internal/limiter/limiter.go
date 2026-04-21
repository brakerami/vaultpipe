// Package limiter provides a concurrency limiter that caps the number of
// simultaneous in-flight secret fetch operations.
package limiter

import (
	"context"
	"errors"
	"fmt"
)

// ErrLimitExceeded is returned when the concurrency cap is already reached and
// the caller's context is cancelled before a slot becomes available.
var ErrLimitExceeded = errors.New("limiter: concurrency limit exceeded")

// Limiter gates concurrent access to a shared resource using a buffered
// semaphore channel.
type Limiter struct {
	sem chan struct{}
	cap int
}

// New creates a Limiter that allows at most n concurrent operations.
// n must be >= 1.
func New(n int) (*Limiter, error) {
	if n < 1 {
		return nil, fmt.Errorf("limiter: capacity must be >= 1, got %d", n)
	}
	return &Limiter{
		sem: make(chan struct{}, n),
		cap: n,
	}, nil
}

// Acquire blocks until a concurrency slot is available or ctx is done.
// On success the caller MUST call Release exactly once.
func (l *Limiter) Acquire(ctx context.Context) error {
	select {
	case l.sem <- struct{}{}:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("%w: %w", ErrLimitExceeded, ctx.Err())
	}
}

// Release frees one concurrency slot.
func (l *Limiter) Release() {
	select {
	case <-l.sem:
	default:
		// no-op: guard against misuse
	}
}

// Cap returns the maximum concurrency configured for this Limiter.
func (l *Limiter) Cap() int {
	return l.cap
}

// InFlight returns the number of slots currently held.
func (l *Limiter) InFlight() int {
	return len(l.sem)
}
