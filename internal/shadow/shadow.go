// Package shadow maintains a secondary copy of resolved secret values
// and detects drift between the live value and the last known good state.
package shadow

import (
	"sync"
	"time"
)

// Entry holds a shadowed secret value alongside the time it was recorded.
type Entry struct {
	Value     string
	RecordedAt time.Time
}

// Shadow stores the last-known values for a set of secret keys.
type Shadow struct {
	mu      sync.RWMutex
	entries map[string]Entry
}

// New returns an empty Shadow store.
func New() *Shadow {
	return &Shadow{entries: make(map[string]Entry)}
}

// Set records the current value for key.
func (s *Shadow) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[key] = Entry{Value: value, RecordedAt: time.Now()}
}

// Get returns the shadowed entry for key and whether it exists.
func (s *Shadow) Get(key string) (Entry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[key]
	return e, ok
}

// Drifted returns true when the live value differs from the shadowed value.
// If no shadow exists for key, Drifted returns false.
func (s *Shadow) Drifted(key, liveValue string) bool {
	e, ok := s.Get(key)
	if !ok {
		return false
	}
	return e.Value != liveValue
}

// Delete removes the shadowed entry for key.
func (s *Shadow) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, key)
}

// Keys returns all keys currently tracked by the shadow store.
func (s *Shadow) Keys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	keys := make([]string, 0, len(s.entries))
	for k := range s.entries {
		keys = append(keys, k)
	}
	return keys
}
