// Package suppress provides a mechanism to suppress duplicate error
// notifications within a configurable window. It prevents alert fatigue
// by ensuring the same error class is only reported once per interval.
package suppress

import (
	"sync"
	"time"
)

// Suppressor tracks error keys and suppresses repeated occurrences
// within a fixed window duration.
type Suppressor struct {
	mu       sync.Mutex
	window   time.Duration
	seen     map[string]time.Time
	nowFn    func() time.Time
}

// New creates a Suppressor with the given deduplication window.
// Panics if window is zero or negative.
func New(window time.Duration) *Suppressor {
	if window <= 0 {
		panic("suppress: window must be positive")
	}
	return &Suppressor{
		window: window,
		seen:   make(map[string]time.Time),
		nowFn:  time.Now,
	}
}

// Allow returns true if the key has not been seen within the current
// window, recording it for future suppression. Returns false if the
// key was already reported recently.
func (s *Suppressor) Allow(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := s.nowFn()
	if last, ok := s.seen[key]; ok && now.Sub(last) < s.window {
		return false
	}
	s.seen[key] = now
	return true
}

// Reset clears all suppressed keys, allowing them to be reported again.
func (s *Suppressor) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.seen = make(map[string]time.Time)
}

// Len returns the number of currently tracked keys.
func (s *Suppressor) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.seen)
}
