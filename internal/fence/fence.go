// Package fence provides a write-once token barrier that prevents
// a secret resolution cycle from being executed more than once
// concurrently. A Fence issues a single-use token; any caller that
// cannot acquire the token is told to wait or back off.
package fence

import (
	"context"
	"errors"
	"sync"
)

// ErrFenced is returned when a caller attempts to acquire a token
// that is already held by another goroutine.
var ErrFenced = errors.New("fence: operation already in progress")

// Fence serialises a single logical operation identified by a string key.
type Fence struct {
	mu    sync.Mutex
	gates map[string]chan struct{}
}

// New returns an initialised Fence.
func New() *Fence {
	return &Fence{gates: make(map[string]chan struct{})}
}

// Acquire attempts to take ownership of key. On success it returns a
// release function that MUST be called when the operation is complete.
// If the key is already held ErrFenced is returned immediately.
func (f *Fence) Acquire(ctx context.Context, key string) (release func(), err error) {
	f.mu.Lock()
	if _, ok := f.gates[key]; ok {
		f.mu.Unlock()
		return nil, ErrFenced
	}
	ch := make(chan struct{})
	f.gates[key] = ch
	f.mu.Unlock()

	release = func() {
		f.mu.Lock()
		delete(f.gates, key)
		f.mu.Unlock()
		close(ch)
	}
	return release, nil
}

// Wait blocks until the operation identified by key has finished, or
// ctx is cancelled. It returns nil if the gate was released normally.
func (f *Fence) Wait(ctx context.Context, key string) error {
	f.mu.Lock()
	ch, ok := f.gates[key]
	f.mu.Unlock()
	if !ok {
		return nil
	}
	select {
	case <-ch:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Active returns the number of keys currently held.
func (f *Fence) Active() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.gates)
}
