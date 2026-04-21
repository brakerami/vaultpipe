// Package snapshot captures and restores the resolved secret environment
// at a point in time, enabling rollback and diff on rotation.
package snapshot

import (
	"sync"
	"time"
)

// Entry holds a single resolved secret value and the time it was captured.
type Entry struct {
	Key       string
	Value     string
	CapturedAt time.Time
}

// Snapshot is an immutable point-in-time view of resolved secrets.
type Snapshot struct {
	mu      sync.RWMutex
	entries map[string]Entry
	takenAt time.Time
}

// Take creates a new Snapshot from the provided key-value pairs.
func Take(secrets map[string]string) *Snapshot {
	now := time.Now().UTC()
	entries := make(map[string]Entry, len(secrets))
	for k, v := range secrets {
		entries[k] = Entry{Key: k, Value: v, CapturedAt: now}
	}
	return &Snapshot{entries: entries, takenAt: now}
}

// TakenAt returns the time the snapshot was created.
func (s *Snapshot) TakenAt() time.Time {
	return s.takenAt
}

// Get returns the value for key and whether it existed.
func (s *Snapshot) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[key]
	return e.Value, ok
}

// Keys returns all keys in the snapshot.
func (s *Snapshot) Keys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	keys := make([]string, 0, len(s.entries))
	for k := range s.entries {
		keys = append(keys, k)
	}
	return keys
}

// ToMap returns a copy of the snapshot as a plain map.
func (s *Snapshot) ToMap() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]string, len(s.entries))
	for k, e := range s.entries {
		out[k] = e.Value
	}
	return out
}
