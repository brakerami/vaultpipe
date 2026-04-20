// Package watch provides secret lease watching and renewal notifications.
// When a secret's TTL is close to expiry, the watcher triggers a callback
// so the caller can refresh the value before it becomes invalid.
package watch

import (
	"context"
	"sync"
	"time"
)

// RenewFunc is called when a watched secret needs renewal.
type RenewFunc func(ctx context.Context, ref string) error

// Watcher monitors secret leases and triggers renewal before expiry.
type Watcher struct {
	mu      sync.Mutex
	entries map[string]*entry
	renew   RenewFunc
	thresh  float64 // fraction of TTL remaining that triggers renewal
}

type entry struct {
	ref     string
	expires time.Time
	ttl     time.Duration
	cancel  context.CancelFunc
}

// New creates a Watcher that calls renew when a secret has less than
// thresholdFraction of its TTL remaining (e.g. 0.25 = 25%).
func New(renew RenewFunc, thresholdFraction float64) *Watcher {
	if thresholdFraction <= 0 || thresholdFraction >= 1 {
		thresholdFraction = 0.25
	}
	return &Watcher{
		entries: make(map[string]*entry),
		renew:   renew,
		thresh:  thresholdFraction,
	}
}

// Add registers a secret ref for watching with the given TTL.
// If the ref is already watched, its lease is reset.
func (w *Watcher) Add(ctx context.Context, ref string, ttl time.Duration) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if e, ok := w.entries[ref]; ok {
		e.cancel()
	}

	watchCtx, cancel := context.WithCancel(ctx)
	e := &entry{
		ref:     ref,
		expires: time.Now().Add(ttl),
		ttl:     ttl,
		cancel:  cancel,
	}
	w.entries[ref] = e
	go w.loop(watchCtx, e)
}

// Remove stops watching the given ref.
func (w *Watcher) Remove(ref string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if e, ok := w.entries[ref]; ok {
		e.cancel()
		delete(w.entries, ref)
	}
}

// Stop cancels all active watchers.
func (w *Watcher) Stop() {
	w.mu.Lock()
	defer w.mu.Unlock()
	for ref, e := range w.entries {
		e.cancel()
		delete(w.entries, ref)
	}
}

func (w *Watcher) loop(ctx context.Context, e *entry) {
	renewalAt := e.expires.Add(-time.Duration(float64(e.ttl) * w.thresh))
	waitFor := time.Until(renewalAt)
	if waitFor < 0 {
		waitFor = 0
	}

	select {
	case <-ctx.Done():
		return
	case <-time.After(waitFor):
		_ = w.renew(ctx, e.ref)
	}
}
