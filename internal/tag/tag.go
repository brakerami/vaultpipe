// Package tag provides utilities for attaching and querying metadata labels
// on secret references, enabling filtering and grouping during resolution.
package tag

import (
	"fmt"
	"strings"
)

// Set holds an ordered collection of key=value tag pairs.
type Set struct {
	tags map[string]string
	keys []string // preserves insertion order
}

// New returns an empty Set.
func New() *Set {
	return &Set{tags: make(map[string]string)}
}

// Add registers a tag. key and value must be non-empty and key must not
// contain '=' or ','.
func (s *Set) Add(key, value string) error {
	if key == "" {
		return fmt.Errorf("tag: key must not be empty")
	}
	if value == "" {
		return fmt.Errorf("tag: value must not be empty for key %q", key)
	}
	if strings.ContainsAny(key, "=,") {
		return fmt.Errorf("tag: key %q contains reserved character", key)
	}
	if _, exists := s.tags[key]; !exists {
		s.keys = append(s.keys, key)
	}
	s.tags[key] = value
	return nil
}

// Get returns the value for key and whether it was present.
func (s *Set) Get(key string) (string, bool) {
	v, ok := s.tags[key]
	return v, ok
}

// Has reports whether key exists in the set.
func (s *Set) Has(key string) bool {
	_, ok := s.tags[key]
	return ok
}

// Len returns the number of tags.
func (s *Set) Len() int { return len(s.tags) }

// String serialises the set as "k1=v1,k2=v2" in insertion order.
func (s *Set) String() string {
	parts := make([]string, 0, len(s.keys))
	for _, k := range s.keys {
		parts = append(parts, k+"="+s.tags[k])
	}
	return strings.Join(parts, ",")
}

// Parse deserialises a string produced by String back into a Set.
// Returns an error if any pair is malformed.
func Parse(raw string) (*Set, error) {
	s := New()
	if raw == "" {
		return s, nil
	}
	for _, pair := range strings.Split(raw, ",") {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("tag: malformed pair %q", pair)
		}
		if err := s.Add(parts[0], parts[1]); err != nil {
			return nil, err
		}
	}
	return s, nil
}
