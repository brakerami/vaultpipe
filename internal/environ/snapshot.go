// Package environ provides utilities for capturing and restoring
// process environment snapshots used during secret injection.
package environ

import (
	"os"
	"strings"
)

// Snapshot holds a point-in-time copy of environment variables.
type Snapshot struct {
	vars map[string]string
}

// Capture returns a Snapshot of the current process environment.
func Capture() *Snapshot {
	raw := os.Environ()
	vars := make(map[string]string, len(raw))
	for _, entry := range raw {
		key, value, found := strings.Cut(entry, "=")
		if !found {
			continue
		}
		vars[key] = value
	}
	return &Snapshot{vars: vars}
}

// FromMap builds a Snapshot from an explicit key/value map.
func FromMap(m map[string]string) *Snapshot {
	copy := make(map[string]string, len(m))
	for k, v := range m {
		copy[k] = v
	}
	return &Snapshot{vars: copy}
}

// Get returns the value for key and whether it was present.
func (s *Snapshot) Get(key string) (string, bool) {
	v, ok := s.vars[key]
	return v, ok
}

// Keys returns all environment variable names in the snapshot.
func (s *Snapshot) Keys() []string {
	keys := make([]string, 0, len(s.vars))
	for k := range s.vars {
		keys = append(keys, k)
	}
	return keys
}

// Environ returns the snapshot as a slice of "KEY=VALUE" strings
// suitable for use with exec.Cmd.Env.
func (s *Snapshot) Environ() []string {
	out := make([]string, 0, len(s.vars))
	for k, v := range s.vars {
		out = append(out, k+"="+v)
	}
	return out
}

// Merge returns a new Snapshot with the provided overrides applied on
// top of the receiver. The original snapshot is not modified.
func (s *Snapshot) Merge(overrides map[string]string) *Snapshot {
	merged := make(map[string]string, len(s.vars)+len(overrides))
	for k, v := range s.vars {
		merged[k] = v
	}
	for k, v := range overrides {
		merged[k] = v
	}
	return &Snapshot{vars: merged}
}

// Len returns the number of variables in the snapshot.
func (s *Snapshot) Len() int {
	return len(s.vars)
}
