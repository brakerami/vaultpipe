// Package drain provides utilities for gracefully draining in-flight
// secret reads and process teardown before vaultpipe exits.
package drain

import (
	"context"
	"sync"
	"time"
)

// DefaultTimeout is the maximum time Drain will wait for work to finish.
const DefaultTimeout = 5 * time.Second

// Drainer tracks active work units and blocks shutdown until they complete
// or the deadline is exceeded.
type Drainer struct {
	mu      sync.Mutex
	wg      sync.WaitGroup
	closed  bool
	timeout time.Duration
}

// New returns a Drainer with the given shutdown timeout.
// If timeout is zero, DefaultTimeout is used.
func New(timeout time.Duration) *Drainer {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	return &Drainer{timeout: timeout}
}

// Acquire marks one unit of work as started.
// It returns false if the Drainer has already been drained (closed).
func (d *Drainer) Acquire() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.closed {
		return false
	}
	d.wg.Add(1)
	return true
}

// Release marks one unit of work as finished.
func (d *Drainer) Release() {
	d.wg.Done()
}

// Drain waits for all acquired work to finish or for ctx / the internal
// timeout to expire, whichever comes first.
// It marks the Drainer as closed so no new work can be acquired after
// Drain returns.
func (d *Drainer) Drain(ctx context.Context) error {
	d.mu.Lock()
	d.closed = true
	d.mu.Unlock()

	deadline, cancel := context.WithTimeout(ctx, d.timeout)
	defer cancel()

	done := make(chan struct{})
	go func() {
		d.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-deadline.Done():
		return deadline.Err()
	}
}
