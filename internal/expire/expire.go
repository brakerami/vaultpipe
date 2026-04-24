// Package expire provides utilities for tracking and reacting to
// secret expiration deadlines relative to a configurable threshold.
package expire

import (
	"sync"
	"time"
)

// Handler is called when a tracked entry crosses the warning threshold.
type Handler func(key string, remaining time.Duration)

// Tracker monitors registered keys and fires a Handler when the time
// remaining before expiry falls at or below the warning threshold.
type Tracker struct {
	mu        sync.Mutex
	entries   map[string]time.Time
	threshold time.Duration
	handler   Handler
}

// New returns a Tracker that will invoke h when a key's remaining
// lifetime is at or below threshold. threshold must be positive.
func New(threshold time.Duration, h Handler) *Tracker {
	if threshold <= 0 {
		panic("expire: threshold must be positive")
	}
	if h == nil {
		panic("expire: handler must not be nil")
	}
	return &Tracker{
		entries:   make(map[string]time.Time),
		threshold: threshold,
		handler:   h,
	}
}

// Track registers key with the given expiry time, replacing any
// existing entry for the same key.
func (t *Tracker) Track(key string, expiresAt time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries[key] = expiresAt
}

// Remove stops tracking key.
func (t *Tracker) Remove(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.entries, key)
}

// Check evaluates all tracked entries against now and invokes the
// Handler for any whose remaining lifetime is at or below the
// configured threshold. Expired entries (remaining <= 0) are also
// reported and then removed.
func (t *Tracker) Check(now time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()

	for key, exp := range t.entries {
		remaining := exp.Sub(now)
		if remaining <= t.threshold {
			t.handler(key, remaining)
			if remaining <= 0 {
				delete(t.entries, key)
			}
		}
	}
}

// Len returns the number of currently tracked entries.
func (t *Tracker) Len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.entries)
}
