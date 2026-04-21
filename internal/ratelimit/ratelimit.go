// Package ratelimit provides a token-bucket rate limiter for Vault API calls.
// It prevents thundering-herd issues when many secrets are resolved at startup.
package ratelimit

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Limiter enforces a maximum number of requests per second using a token bucket.
type Limiter struct {
	mu       sync.Mutex
	tokens   float64
	max      float64
	rate     float64 // tokens per nanosecond
	lastTick time.Time
}

// New creates a Limiter that allows up to rps requests per second.
// rps must be greater than zero.
func New(rps float64) (*Limiter, error) {
	if rps <= 0 {
		return nil, fmt.Errorf("ratelimit: rps must be greater than zero, got %v", rps)
	}
	return &Limiter{
		tokens:   rps,
		max:      rps,
		rate:     rps / float64(time.Second),
		lastTick: time.Now(),
	}, nil
}

// Wait blocks until a token is available or ctx is cancelled.
// Returns ctx.Err() if the context expires before a token is acquired.
func (l *Limiter) Wait(ctx context.Context) error {
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		l.mu.Lock()
		now := time.Now()
		elapsed := now.Sub(l.lastTick)
		l.tokens += float64(elapsed) * l.rate
		if l.tokens > l.max {
			l.tokens = l.max
		}
		l.lastTick = now
		if l.tokens >= 1.0 {
			l.tokens -= 1.0
			l.mu.Unlock()
			return nil
		}
		// Calculate how long until the next token is available.
		wait := time.Duration((1.0-l.tokens)/l.rate)
		l.mu.Unlock()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(wait):
		}
	}
}

// Available returns the current number of available tokens (for observability).
func (l *Limiter) Available() float64 {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.tokens
}
