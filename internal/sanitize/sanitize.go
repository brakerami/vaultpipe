// Package sanitize provides utilities for stripping control characters
// and non-printable bytes from secret values before they are injected
// into a process environment. This prevents terminal escape injection
// and other unexpected behaviour caused by malformed Vault responses.
package sanitize

import (
	"strings"
	"unicode"
)

// Value removes non-printable and control characters from s, returning
// a clean string safe for use as an environment variable value.
// Printable Unicode letters, digits, punctuation, symbols, and ASCII
// spaces are preserved; everything else is dropped.
func Value(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if isAllowed(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// Map applies Value to every entry in m, returning a new map with the
// sanitized values. The original map is not modified.
func Map(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = Value(v)
	}
	return out
}

// isAllowed reports whether r is safe to include in an env var value.
func isAllowed(r rune) bool {
	if r == ' ' || r == '\t' {
		return true
	}
	if unicode.IsControl(r) {
		return false
	}
	if !unicode.IsPrint(r) {
		return false
	}
	return true
}
