// Package latch provides a one-shot boolean gate that can be set exactly once.
// Once latched, subsequent Set calls are no-ops and callers can wait for the
// latch to close via Wait.
package latch

import (
	"context"
	"sync"
)

// Latch is a one-shot gate. The zero value is ready to use.
type Latch struct {
	once sync.Once
	ch   chan struct{}
	mu   sync.Mutex
}

// New returns an initialised Latch.
func New() *Latch {
	return &Latch{ch: make(chan struct{})}
}

func (l *Latch) init() chan struct{} {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.ch == nil {
		l.ch = make(chan struct{})
	}
	return l.ch
}

// Set closes the latch. Subsequent calls are silently ignored.
func (l *Latch) Set() {
	ch := l.init()
	l.once.Do(func() { close(ch) })
}

// IsSet reports whether the latch has been closed.
func (l *Latch) IsSet() bool {
	select {
	case <-l.init():
		return true
	default:
		return false
	}
}

// Wait blocks until the latch is set or ctx is cancelled.
// Returns ctx.Err() if the context expires first, nil otherwise.
func (l *Latch) Wait(ctx context.Context) error {
	select {
	case <-l.init():
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// C returns a channel that is closed when the latch is set.
// Callers must not close the returned channel.
func (l *Latch) C() <-chan struct{} {
	return l.init()
}
