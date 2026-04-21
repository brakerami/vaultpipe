// Package rotate provides secret rotation detection and callback triggering
// for secrets fetched from Vault. It compares newly resolved values against
// previously known ones and fires registered hooks when a change is detected.
package rotate

import (
	"context"
	"sync"
)

// ChangeFunc is called when a secret value changes. oldVal may be empty
// if this is the first observation of the secret.
type ChangeFunc func(ctx context.Context, key, oldVal, newVal string)

// Detector tracks secret values and detects rotation.
type Detector struct {
	mu      sync.Mutex
	known   map[string]string
	hooks   []ChangeFunc
}

// New returns an initialised Detector.
func New() *Detector {
	return &Detector{
		known: make(map[string]string),
	}
}

// OnChange registers a hook that will be called whenever a secret rotates.
func (d *Detector) OnChange(fn ChangeFunc) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.hooks = append(d.hooks, fn)
}

// Observe records the current value for key. If the value differs from the
// previously recorded value all registered hooks are invoked synchronously.
// Returns true if a rotation was detected.
func (d *Detector) Observe(ctx context.Context, key, value string) bool {
	d.mu.Lock()
	old, exists := d.known[key]
	d.known[key] = value
	hooks := make([]ChangeFunc, len(d.hooks))
	copy(hooks, d.hooks)
	d.mu.Unlock()

	if exists && old == value {
		return false
	}
	if !exists {
		// First observation — not a rotation, just seeding.
		return false
	}
	for _, fn := range hooks {
		fn(ctx, key, old, value)
	}
	return true
}

// Seed stores an initial value for key without triggering hooks. Useful for
// pre-populating known values on startup.
func (d *Detector) Seed(key, value string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.known[key] = value
}

// Reset removes all tracked state for key.
func (d *Detector) Reset(key string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.known, key)
}
