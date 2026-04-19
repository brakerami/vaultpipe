// Package mask provides utilities for redacting secret values in output.
package mask

import (
	"strings"
	"sync"
)

const redacted = "[REDACTED]"

// Masker holds a set of known secret values and can redact them from strings.
type Masker struct {
	mu      sync.RWMutex
	secrets map[string]struct{}
}

// New returns an empty Masker.
func New() *Masker {
	return &Masker{
		secrets: make(map[string]struct{}),
	}
}

// Add registers a secret value to be masked.
// Empty strings are ignored.
func (m *Masker) Add(secret string) {
	if secret == "" {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.secrets[secret] = struct{}{}
}

// Mask replaces all registered secret values in s with [REDACTED].
func (m *Masker) Mask(s string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for secret := range m.secrets {
		s = strings.ReplaceAll(s, secret, redacted)
	}
	return s
}

// Len returns the number of registered secrets.
func (m *Masker) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.secrets)
}
