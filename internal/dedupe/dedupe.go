// Package dedupe provides idempotency tracking for secret fetch operations.
// It suppresses duplicate fetches for the same path within a configurable window,
// returning the cached result to all concurrent callers.
package dedupe

import (
	"context"
	"sync"
	"time"
)

// Result holds the outcome of a single fetch.
type Result struct {
	Value string
	Err   error
}

// FetchFunc is the underlying fetch operation to deduplicate.
type FetchFunc func(ctx context.Context, path string) (string, error)

// flight represents an in-progress or recently completed fetch.
type flight struct {
	mu      sync.Mutex
	done    chan struct{}
	result  Result
	expires time.Time
}

// Group deduplicates concurrent and near-concurrent fetches for the same key.
type Group struct {
	mu      sync.Mutex
	window  time.Duration
	flights map[string]*flight
}

// New returns a Group that suppresses duplicate fetches within window.
// A zero window disables result reuse across sequential calls but still
// collapses concurrent callers into a single in-flight request.
func New(window time.Duration) *Group {
	return &Group{
		window:  window,
		flights: make(map[string]*flight),
	}
}

// Do executes fetch for path, collapsing concurrent callers into one request.
// If a successful result is still within the reuse window it is returned
// immediately without calling fetch again.
func (g *Group) Do(ctx context.Context, path string, fetch FetchFunc) (string, error) {
	g.mu.Lock()
	f, ok := g.flights[path]
	if ok {
		f.mu.Lock()
		if f.expires.After(time.Now()) {
			res := f.result
			f.mu.Unlock()
			g.mu.Unlock()
			return res.Value, res.Err
		}
		f.mu.Unlock()
	}

	f = &flight{done: make(chan struct{})}
	g.flights[path] = f
	g.mu.Unlock()

	go func() {
		v, err := fetch(ctx, path)
		f.mu.Lock()
		f.result = Result{Value: v, Err: err}
		if err == nil && g.window > 0 {
			f.expires = time.Now().Add(g.window)
		}
		f.mu.Unlock()
		close(f.done)
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-f.done:
		f.mu.Lock()
		res := f.result
		f.mu.Unlock()
		return res.Value, res.Err
	}
}

// Invalidate removes any cached result for path, forcing the next call to
// execute a fresh fetch.
func (g *Group) Invalidate(path string) {
	g.mu.Lock()
	delete(g.flights, path)
	g.mu.Unlock()
}
