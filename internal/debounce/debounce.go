// Package debounce provides a mechanism to coalesce rapid successive calls
// into a single execution after a quiet period has elapsed. This is useful
// for batching secret-rotation events or config-reload triggers that may
// fire in quick succession.
package debounce

import (
	"context"
	"sync"
	"time"
)

// Func is the signature of the function that will be debounced.
type Func func(ctx context.Context)

// Debouncer delays execution of fn until no further calls have been made
// within the configured wait duration.
type Debouncer struct {
	wait  time.Duration
	fn    Func
	mu    sync.Mutex
	timer *time.Timer
}

// New returns a Debouncer that will invoke fn after wait has elapsed since
// the last call to Trigger. wait must be positive.
func New(wait time.Duration, fn Func) *Debouncer {
	if wait <= 0 {
		panic("debounce: wait must be positive")
	}
	if fn == nil {
		panic("debounce: fn must not be nil")
	}
	return &Debouncer{wait: wait, fn: fn}
}

// Trigger schedules fn to run after the debounce window. If Trigger is called
// again before the window expires the timer is reset.
func (d *Debouncer) Trigger(ctx context.Context) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
	}
	d.timer = time.AfterFunc(d.wait, func() {
		d.fn(ctx)
	})
}

// Flush cancels any pending timer and invokes fn immediately. If no call is
// pending, Flush is a no-op.
func (d *Debouncer) Flush(ctx context.Context) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer == nil {
		return
	}
	if d.timer.Stop() {
		d.timer = nil
		go d.fn(ctx)
	}
}

// Stop cancels any pending invocation without executing fn.
func (d *Debouncer) Stop() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
		d.timer = nil
	}
}
