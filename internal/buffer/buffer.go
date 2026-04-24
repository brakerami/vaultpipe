// Package buffer provides a fixed-capacity ring buffer for collecting
// recent log lines or secret access events, suitable for diagnostics
// and tail-style inspection without unbounded memory growth.
package buffer

import (
	"sync"
	"time"
)

// Entry holds a single buffered record.
type Entry struct {
	At      time.Time
	Message string
}

// Ring is a thread-safe, fixed-capacity circular buffer of Entry values.
// When the buffer is full, the oldest entry is overwritten.
type Ring struct {
	mu       sync.Mutex
	entries  []Entry
	cap      int
	head     int // next write position
	count    int
}

// New returns a Ring with the given capacity. It panics if cap < 1.
func New(cap int) *Ring {
	if cap < 1 {
		panic("buffer: capacity must be at least 1")
	}
	return &Ring{
		entries: make([]Entry, cap),
		cap:     cap,
	}
}

// Add appends a message to the ring, overwriting the oldest entry if full.
func (r *Ring) Add(msg string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[r.head] = Entry{At: time.Now(), Message: msg}
	r.head = (r.head + 1) % r.cap
	if r.count < r.cap {
		r.count++
	}
}

// Snapshot returns a copy of all buffered entries in chronological order.
func (r *Ring) Snapshot() []Entry {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.count == 0 {
		return nil
	}
	out := make([]Entry, r.count)
	start := (r.head - r.count + r.cap) % r.cap
	for i := 0; i < r.count; i++ {
		out[i] = r.entries[(start+i)%r.cap]
	}
	return out
}

// Len returns the number of entries currently held.
func (r *Ring) Len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.count
}

// Reset discards all buffered entries.
func (r *Ring) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.head = 0
	r.count = 0
}
