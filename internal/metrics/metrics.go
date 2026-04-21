// Package metrics provides lightweight counters and gauges for tracking
// vaultpipe runtime events such as secret fetches, cache hits, and errors.
package metrics

import (
	"sync"
	"sync/atomic"
)

// Counter is a monotonically increasing integer counter.
type Counter struct {
	value uint64
}

// Inc increments the counter by 1.
func (c *Counter) Inc() {
	atomic.AddUint64(&c.value, 1)
}

// Value returns the current counter value.
func (c *Counter) Value() uint64 {
	return atomic.LoadUint64(&c.value)
}

// Registry holds named counters for runtime instrumentation.
type Registry struct {
	mu       sync.RWMutex
	counters map[string]*Counter
}

// New creates an empty Registry.
func New() *Registry {
	return &Registry{
		counters: make(map[string]*Counter),
	}
}

// Counter returns the named counter, creating it if it does not exist.
func (r *Registry) Counter(name string) *Counter {
	r.mu.RLock()
	c, ok := r.counters[name]
	r.mu.RUnlock()
	if ok {
		return c
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	// Double-checked locking.
	if c, ok = r.counters[name]; ok {
		return c
	}
	c = &Counter{}
	r.counters[name] = c
	return c
}

// Snapshot returns a point-in-time copy of all counter values keyed by name.
func (r *Registry) Snapshot() map[string]uint64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[string]uint64, len(r.counters))
	for name, c := range r.counters {
		out[name] = c.Value()
	}
	return out
}

// Reset sets all counters back to zero. Primarily useful in tests.
func (r *Registry) Reset() {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, c := range r.counters {
		atomic.StoreUint64(&c.value, 0)
	}
}
