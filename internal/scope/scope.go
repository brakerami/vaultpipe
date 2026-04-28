// Package scope provides key-prefix scoping for secret lookups,
// allowing callers to namespace environment variable keys by a
// configurable prefix string.
package scope

import (
	"errors"
	"strings"
)

// Scope applies a consistent prefix to secret keys before they are
// injected into the process environment.
type Scope struct {
	prefix    string
	separator string
}

// New returns a Scope with the given prefix. The separator defaults
// to "_" when empty. An error is returned if the prefix contains
// characters that are invalid in environment variable names.
func New(prefix, separator string) (*Scope, error) {
	if prefix == "" {
		return nil, errors.New("scope: prefix must not be empty")
	}
	if strings.ContainsAny(prefix, "= \t\n") {
		return nil, errors.New("scope: prefix contains invalid characters")
	}
	if separator == "" {
		separator = "_"
	}
	return &Scope{prefix: strings.ToUpper(prefix), separator: separator}, nil
}

// Apply prepends the scope prefix to key, returning the scoped key.
// key is upper-cased before the prefix is applied.
func (s *Scope) Apply(key string) string {
	if key == "" {
		return s.prefix
	}
	return s.prefix + s.separator + strings.ToUpper(key)
}

// Strip removes the scope prefix from key. If key does not carry
// the prefix, the original value is returned unchanged along with
// false to indicate no match.
func (s *Scope) Strip(key string) (string, bool) {
	upper := strings.ToUpper(key)
	expected := s.prefix + s.separator
	if strings.HasPrefix(upper, expected) {
		return upper[len(expected):], true
	}
	return key, false
}

// Prefix returns the configured prefix string.
func (s *Scope) Prefix() string { return s.prefix }

// Map applies the scope prefix to every key in src, returning a new
// map with the scoped keys and the original values.
func (s *Scope) Map(src map[string]string) map[string]string {
	out := make(map[string]string, len(src))
	for k, v := range src {
		out[s.Apply(k)] = v
	}
	return out
}
