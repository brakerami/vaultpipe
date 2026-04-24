// Package cooldown enforces a minimum interval between successive
// executions of the same keyed operation. It is useful for preventing
// burst re-fetches of secrets after a rotation event.
package cooldown

import (
	"sync"
	"time"
)

// Cooldown tracks the last execution time per key and reports whether
// enough time has elapsed to allow the next execution.
type Cooldown struct {
	mu       sync.Mutex
	interval time.Duration
	last     map[string]time.Time
	now      func() time.Time
}

// New returns a Cooldown that enforces the given minimum interval between
// successive calls for the same key. It panics if interval is <= 0.
func New(interval time.Duration) *Cooldown {
	if interval <= 0 {
		panic("cooldown: interval must be positive")
	}
	return &Cooldown{
		interval: interval,
		last:     make(map[string]time.Time),
		now:      time.Now,
	}
}

// Allow reports whether the operation identified by key may proceed.
// If the cooldown period has elapsed (or the key has never been seen),
// Allow records the current time and returns true. Otherwise it returns
// false without updating the recorded time.
func (c *Cooldown) Allow(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := c.now()
	if t, ok := c.last[key]; ok && now.Sub(t) < c.interval {
		return false
	}
	c.last[key] = now
	return true
}

// Reset removes the recorded timestamp for key, allowing the next call
// to Allow to succeed immediately.
func (c *Cooldown) Reset(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.last, key)
}

// Remaining returns the duration that must still elapse before Allow
// will return true for key. It returns 0 if the cooldown has expired or
// the key is unknown.
func (c *Cooldown) Remaining(key string) time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()

	t, ok := c.last[key]
	if !ok {
		return 0
	}
	elapsed := c.now().Sub(t)
	if elapsed >= c.interval {
		return 0
	}
	return c.interval - elapsed
}
