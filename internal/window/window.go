// Package window implements a sliding-window counter used to track
// event rates over a rolling time interval.
package window

import (
	"sync"
	"time"
)

// Window is a thread-safe sliding-window counter.
type Window struct {
	mu       sync.Mutex
	size     time.Duration
	buckets  int
	counts   []int64
	times    []time.Time
	cursor   int
	now      func() time.Time
}

// New returns a Window that tracks events over the given duration split
// into the requested number of buckets. buckets must be >= 1.
func New(size time.Duration, buckets int) *Window {
	if buckets < 1 {
		buckets = 1
	}
	return &Window{
		size:    size,
		buckets: buckets,
		counts:  make([]int64, buckets),
		times:   make([]time.Time, buckets),
		now:     time.Now,
	}
}

// Add records n events at the current time.
func (w *Window) Add(n int64) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.advance()
	w.counts[w.cursor] += n
}

// Count returns the total number of events within the sliding window.
func (w *Window) Count() int64 {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.advance()
	cutoff := w.now().Add(-w.size)
	var total int64
	for i, t := range w.times {
		if !t.IsZero() && t.After(cutoff) {
			total += w.counts[i]
		}
	}
	return total
}

// Reset clears all recorded events.
func (w *Window) Reset() {
	w.mu.Lock()
	defer w.mu.Unlock()
	for i := range w.counts {
		w.counts[i] = 0
		w.times[i] = time.Time{}
	}
}

// advance moves the cursor to the current bucket, clearing stale ones.
func (w *Window) advance() {
	now := w.now()
	bucketDur := w.size / time.Duration(w.buckets)
	next := (w.cursor + 1) % w.buckets
	if w.times[w.cursor].IsZero() {
		w.times[w.cursor] = now
		return
	}
	if now.Sub(w.times[w.cursor]) >= bucketDur {
		w.cursor = next
		w.counts[w.cursor] = 0
		w.times[w.cursor] = now
	}
}
