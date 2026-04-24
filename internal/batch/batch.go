// Package batch provides concurrent secret resolution with bounded parallelism.
package batch

import (
	"context"
	"sync"
)

// Result holds the outcome of resolving a single secret reference.
type Result struct {
	Key   string
	Value string
	Err   error
}

// FetchFunc resolves a secret path to its plaintext value.
type FetchFunc func(ctx context.Context, path string) (string, error)

// Resolver maps environment variable keys to secret paths and resolves them
// concurrently, honouring a maximum parallelism cap.
type Resolver struct {
	fetch       FetchFunc
	concurrency int
}

// New returns a Resolver that will call fetch for each secret, running at most
// concurrency goroutines simultaneously. concurrency must be >= 1.
func New(fetch FetchFunc, concurrency int) *Resolver {
	if concurrency < 1 {
		concurrency = 1
	}
	return &Resolver{fetch: fetch, concurrency: concurrency}
}

// Resolve resolves all entries in refs (key -> path) concurrently.
// It returns one Result per entry. The first context cancellation or deadline
// causes in-flight fetches to abort; already-started goroutines drain before
// Resolve returns.
func (r *Resolver) Resolve(ctx context.Context, refs map[string]string) []Result {
	type work struct {
		key  string
		path string
	}

	jobs := make(chan work, len(refs))
	for k, p := range refs {
		jobs <- work{key: k, path: p}
	}
	close(jobs)

	results := make([]Result, 0, len(refs))
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < r.concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				val, err := r.fetch(ctx, j.path)
				mu.Lock()
				results = append(results, Result{Key: j.key, Value: val, Err: err})
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	return results
}
