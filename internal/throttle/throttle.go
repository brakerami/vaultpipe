// Package throttle provides a time-based request throttler that limits
// how frequently an operation can be triggered, regardless of how many
// callers request it. Unlike ratelimit (token bucket), throttle enforces
// a minimum interval between successive executions.
package throttle

import (
	"context"
	"sync"
	"time"
)

// Throttle enforces a minimum interval between successive calls to Do.
type Throttle struct {
	mu       sync.Mutex
	interval time.Duration
	lastRun  time.Time
}

// New creates a Throttle that allows at most one execution per interval.
// interval must be positive; otherwise New panics.
func New(interval time.Duration) *Throttle {
	if interval <= 0 {
		panic("throttle: interval must be positive")
	}
	return &Throttle{interval: interval}
}

// Do calls fn if at least one full interval has elapsed since the last
// successful call. If the throttle is still within the cooldown window,
// Do returns immediately without calling fn.
//
// Do respects ctx cancellation while waiting for the remaining cooldown.
// If ctx is cancelled during the wait, Do returns ctx.Err().
func (t *Throttle) Do(ctx context.Context, fn func() error) error {
	t.mu.Lock()
	now := time.Now()
	next := t.lastRun.Add(t.interval)
	wait := next.Sub(now)
	if wait <= 0 {
		t.lastRun = now
		t.mu.Unlock()
		return fn()
	}
	t.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(wait):
	}

	t.mu.Lock()
	t.lastRun = time.Now()
	t.mu.Unlock()
	return fn()
}

// Reset clears the last-run timestamp, allowing the next call to Do to
// execute immediately regardless of when it last ran.
func (t *Throttle) Reset() {
	t.mu.Lock()
	t.lastRun = time.Time{}
	t.mu.Unlock()
}

// LastRun returns the time of the most recent successful execution, or
// the zero value if Do has never been called.
func (t *Throttle) LastRun() time.Time {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.lastRun
}
