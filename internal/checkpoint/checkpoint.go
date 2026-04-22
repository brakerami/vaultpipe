// Package checkpoint provides a simple mechanism for tracking and persisting
// named progress markers across pipeline stages. Each checkpoint records
// whether a named stage has been reached, along with a timestamp and optional
// metadata. This is useful for resumable operations and observability.
package checkpoint

import (
	"fmt"
	"sync"
	"time"
)

// Status represents the completion state of a checkpoint.
type Status string

const (
	// StatusPending indicates the checkpoint has not yet been reached.
	StatusPending Status = "pending"
	// StatusReached indicates the checkpoint was successfully passed.
	StatusReached Status = "reached"
	// StatusFailed indicates the checkpoint was attempted but failed.
	StatusFailed Status = "failed"
)

// Entry records the state of a single named checkpoint.
type Entry struct {
	Name      string
	Status    Status
	ReachedAt time.Time
	Err       error
}

// Tracker maintains a set of named checkpoints and their states.
type Tracker struct {
	mu      sync.RWMutex
	entries map[string]*Entry
}

// New returns an initialised Tracker with no checkpoints recorded.
func New() *Tracker {
	return &Tracker{
		entries: make(map[string]*Entry),
	}
}

// Mark records a checkpoint as reached at the current time.
// Calling Mark on an already-reached checkpoint updates the timestamp.
func (t *Tracker) Mark(name string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.entries[name] = &Entry{
		Name:      name,
		Status:    StatusReached,
		ReachedAt: time.Now().UTC(),
	}
}

// Fail records a checkpoint as failed, capturing the associated error.
func (t *Tracker) Fail(name string, err error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.entries[name] = &Entry{
		Name:      name,
		Status:    StatusFailed,
		ReachedAt: time.Now().UTC(),
		Err:       err,
	}
}

// Get returns the Entry for the named checkpoint.
// If the checkpoint has never been recorded, a pending Entry is returned
// along with false to indicate a miss.
func (t *Tracker) Get(name string) (Entry, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	e, ok := t.entries[name]
	if !ok {
		return Entry{Name: name, Status: StatusPending}, false
	}
	return *e, true
}

// Reached reports whether the named checkpoint has been successfully marked.
func (t *Tracker) Reached(name string) bool {
	e, ok := t.Get(name)
	return ok && e.Status == StatusReached
}

// All returns a snapshot of every recorded Entry, keyed by name.
func (t *Tracker) All() map[string]Entry {
	t.mu.RLock()
	defer t.mu.RUnlock()

	out := make(map[string]Entry, len(t.entries))
	for k, v := range t.entries {
		out[k] = *v
	}
	return out
}

// Reset removes all recorded checkpoints, returning the Tracker to its
// initial state.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.entries = make(map[string]*Entry)
}

// Summary returns a human-readable string listing all checkpoint statuses.
func (t *Tracker) Summary() string {
	all := t.All()
	if len(all) == 0 {
		return "no checkpoints recorded"
	}

	var s string
	for name, e := range all {
		if e.Err != nil {
			s += fmt.Sprintf("%s: %s (%v)\n", name, e.Status, e.Err)
		} else {
			s += fmt.Sprintf("%s: %s\n", name, e.Status)
		}
	}
	return s
}
