// Package prefetch provides background pre-loading of secrets before
// their leases expire, reducing latency spikes during rotation.
package prefetch

import (
	"context"
	"sync"
	"time"
)

// FetchFunc retrieves a secret value for the given path.
type FetchFunc func(ctx context.Context, path string) (string, error)

// Entry holds the state for a single prefetch registration.
type Entry struct {
	path      string
	threshold float64 // fraction of TTL remaining that triggers prefetch [0,1)
	ttl       time.Duration
	started   time.Time
	cancel    context.CancelFunc
}

// Prefetcher manages background prefetch goroutines for registered secrets.
type Prefetcher struct {
	mu      sync.Mutex
	entries map[string]*Entry
	fetch   FetchFunc
	onRenew func(path, value string)
}

// New creates a Prefetcher that calls fetch to refresh a secret and invokes
// onRenew with the updated value.
func New(fetch FetchFunc, onRenew func(path, value string)) *Prefetcher {
	if onRenew == nil {
		onRenew = func(_, _ string) {}
	}
	return &Prefetcher{
		entries: make(map[string]*Entry),
		fetch:   fetch,
		onRenew: onRenew,
	}
}

// Register schedules a prefetch for path after (1-threshold)*ttl has elapsed.
// Calling Register for an existing path replaces the previous schedule.
func (p *Prefetcher) Register(path string, ttl time.Duration, threshold float64) {
	if threshold <= 0 || threshold >= 1 {
		threshold = 0.2
	}
	p.mu.Lock()
	if e, ok := p.entries[path]; ok {
		e.cancel()
	}
	ctx, cancel := context.WithCancel(context.Background())
	e := &Entry{
		path:      path,
		threshold: threshold,
		ttl:       ttl,
		started:   time.Now(),
		cancel:    cancel,
	}
	p.entries[path] = e
	p.mu.Unlock()

	delay := time.Duration(float64(ttl) * (1 - threshold))
	go p.run(ctx, path, delay)
}

// Deregister cancels any pending prefetch for path.
func (p *Prefetcher) Deregister(path string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if e, ok := p.entries[path]; ok {
		e.cancel()
		delete(p.entries, path)
	}
}

// Stop cancels all registered prefetches.
func (p *Prefetcher) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()
	for path, e := range p.entries {
		e.cancel()
		delete(p.entries, path)
	}
}

func (p *Prefetcher) run(ctx context.Context, path string, delay time.Duration) {
	select {
	case <-ctx.Done():
		return
	case <-time.After(delay):
	}
	val, err := p.fetch(ctx, path)
	if err != nil || ctx.Err() != nil {
		return
	}
	p.onRenew(path, val)
}
