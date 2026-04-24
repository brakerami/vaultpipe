// Package stagger spreads concurrent operations across a time window
// to avoid thundering-herd effects when many leases expire simultaneously.
package stagger

import (
	"context"
	"math/rand"
	"sync"
	"time"
)

// Stagger holds configuration for spreading work over a window.
type Stagger struct {
	mu     sync.Mutex
	window time.Duration
	rng    *rand.Rand
}

// New returns a Stagger that spreads calls over the given window.
// Panics if window is zero or negative.
func New(window time.Duration) *Stagger {
	if window <= 0 {
		panic("stagger: window must be positive")
	}
	return &Stagger{
		window: window,
		rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Delay returns a random duration in [0, window).
func (s *Stagger) Delay() time.Duration {
	s.mu.Lock()
	d := time.Duration(s.rng.Int63n(int64(s.window)))
	s.mu.Unlock()
	return d
}

// Wait blocks for a random duration within the window, or until ctx is
// cancelled. Returns ctx.Err() if the context is cancelled before the
// delay elapses, nil otherwise.
func (s *Stagger) Wait(ctx context.Context) error {
	delay := s.Delay()
	if delay == 0 {
		return ctx.Err()
	}
	select {
	case <-time.After(delay):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Do waits a staggered delay then calls fn. If the context is cancelled
// during the wait, fn is not called and the context error is returned.
func (s *Stagger) Do(ctx context.Context, fn func() error) error {
	if err := s.Wait(ctx); err != nil {
		return err
	}
	return fn()
}
