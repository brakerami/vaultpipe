// Package quota enforces per-key request quotas over a sliding window,
// preventing any single secret path from being fetched too frequently.
package quota

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// ErrExceeded is returned when a key has exceeded its allowed request quota.
var ErrExceeded = errors.New("quota exceeded")

// Config holds the parameters for a Quota.
type Config struct {
	// MaxRequests is the maximum number of requests allowed per Window.
	MaxRequests int
	// Window is the duration over which MaxRequests is measured.
	Window time.Duration
}

// entry tracks request timestamps for a single key.
type entry struct {
	times []time.Time
}

// Quota enforces a sliding-window rate limit per key.
type Quota struct {
	cfg    Config
	mu     sync.Mutex
	keys   map[string]*entry
}

// New creates a new Quota with the given configuration.
// Returns an error if MaxRequests < 1 or Window <= 0.
func New(cfg Config) (*Quota, error) {
	if cfg.MaxRequests < 1 {
		return nil, fmt.Errorf("quota: MaxRequests must be at least 1, got %d", cfg.MaxRequests)
	}
	if cfg.Window <= 0 {
		return nil, fmt.Errorf("quota: Window must be positive, got %s", cfg.Window)
	}
	return &Quota{
		cfg:  cfg,
		keys: make(map[string]*entry),
	}, nil
}

// Allow checks whether the given key is within quota and records the attempt.
// Returns ErrExceeded if the quota has been reached for the current window.
func (q *Quota) Allow(key string) error {
	now := time.Now()
	cutoff := now.Add(-q.cfg.Window)

	q.mu.Lock()
	defer q.mu.Unlock()

	e, ok := q.keys[key]
	if !ok {
		e = &entry{}
		q.keys[key] = e
	}

	// Evict timestamps outside the window.
	valid := e.times[:0]
	for _, t := range e.times {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}
	e.times = valid

	if len(e.times) >= q.cfg.MaxRequests {
		return fmt.Errorf("%w: key %q reached %d requests in %s", ErrExceeded, key, q.cfg.MaxRequests, q.cfg.Window)
	}

	e.times = append(e.times, now)
	return nil
}

// Reset clears all recorded request history for the given key.
func (q *Quota) Reset(key string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	delete(q.keys, key)
}

// Count returns the number of requests recorded for key within the current window.
func (q *Quota) Count(key string) int {
	cutoff := time.Now().Add(-q.cfg.Window)
	q.mu.Lock()
	defer q.mu.Unlock()
	e, ok := q.keys[key]
	if !ok {
		return 0
	}
	count := 0
	for _, t := range e.times {
		if t.After(cutoff) {
			count++
		}
	}
	return count
}
