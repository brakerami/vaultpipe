// Package sieve filters secret references based on configurable allow/deny rules.
// It is used to restrict which Vault paths a process may request at runtime.
package sieve

import (
	"fmt"
	"strings"
	"sync"
)

// Rule describes a single allow or deny pattern.
type Rule struct {
	Allow bool
	// Prefix is matched against the beginning of a secret path.
	Prefix string
}

// Sieve holds an ordered list of rules evaluated top-to-bottom.
// The first matching rule wins. If no rule matches, the path is allowed.
type Sieve struct {
	mu    sync.RWMutex
	rules []Rule
}

// New returns an empty Sieve. Rules may be added via Allow and Deny.
func New() *Sieve {
	return &Sieve{}
}

// Allow appends an allow rule for paths that start with prefix.
func (s *Sieve) Allow(prefix string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rules = append(s.rules, Rule{Allow: true, Prefix: prefix})
}

// Deny appends a deny rule for paths that start with prefix.
func (s *Sieve) Deny(prefix string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rules = append(s.rules, Rule{Allow: false, Prefix: prefix})
}

// Permitted reports whether path is allowed by the current rule set.
// Rules are evaluated in insertion order; the first match decides.
// If no rule matches, the path is permitted.
func (s *Sieve) Permitted(path string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, r := range s.rules {
		if strings.HasPrefix(path, r.Prefix) {
			return r.Allow
		}
	}
	return true
}

// Check returns nil when path is permitted and a descriptive error otherwise.
func (s *Sieve) Check(path string) error {
	if !s.Permitted(path) {
		return fmt.Errorf("sieve: path %q is denied by policy", path)
	}
	return nil
}

// Rules returns a snapshot of the current rule list.
func (s *Sieve) Rules() []Rule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Rule, len(s.rules))
	copy(out, s.rules)
	return out
}

// Reset removes all rules.
func (s *Sieve) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rules = s.rules[:0]
}
