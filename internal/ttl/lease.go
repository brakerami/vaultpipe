package ttl

import (
	"sync"
	"time"
)

// Lease tracks the expiry of a single secret value.
type Lease struct {
	mu        sync.RWMutex
	createdAt time.Time
	duration  time.Duration
}

// NewLease creates a Lease that expires after d from now.
func NewLease(d time.Duration) *Lease {
	return &Lease{
		createdAt: time.Now(),
		duration:  d,
	}
}

// Remaining returns how much time is left before the lease expires.
// Returns 0 if already expired.
func (l *Lease) Remaining() time.Duration {
	l.mu.RLock()
	defer l.mu.RUnlock()

	expiry := l.createdAt.Add(l.duration)
	remaining := time.Until(expiry)
	if remaining < 0 {
		return 0
	}
	return remaining
}

// Expired reports whether the lease has expired.
func (l *Lease) Expired() bool {
	return l.Remaining() == 0
}

// Renew resets the lease start time to now, extending it by its original duration.
func (l *Lease) Renew() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.createdAt = time.Now()
}

// ExpiresAt returns the absolute time at which the lease expires.
func (l *Lease) ExpiresAt() time.Time {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.createdAt.Add(l.duration)
}
