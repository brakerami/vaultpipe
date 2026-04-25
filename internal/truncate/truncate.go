// Package truncate provides utilities for truncating secret values
// to a maximum byte length before injection, preventing oversized
// environment variable values from breaking process environments.
package truncate

import (
	"fmt"
	"unicode/utf8"
)

// DefaultMaxBytes is the default maximum number of bytes allowed for a single
// environment variable value.
const DefaultMaxBytes = 32 * 1024 // 32 KiB

// ErrTruncated is returned when a value exceeds the configured limit.
type ErrTruncated struct {
	Key      string
	Original int
	Limit    int
}

func (e *ErrTruncated) Error() string {
	return fmt.Sprintf("truncate: value for %q truncated from %d to %d bytes", e.Key, e.Original, e.Limit)
}

// Truncator trims secret values that exceed a configured byte limit.
type Truncator struct {
	maxBytes int
}

// New returns a Truncator with the given byte limit.
// If maxBytes is <= 0, DefaultMaxBytes is used.
func New(maxBytes int) *Truncator {
	if maxBytes <= 0 {
		maxBytes = DefaultMaxBytes
	}
	return &Truncator{maxBytes: maxBytes}
}

// Apply trims value to at most t.maxBytes bytes, preserving valid UTF-8
// rune boundaries. It returns the (possibly truncated) string and a non-nil
// *ErrTruncated if truncation occurred.
func (t *Truncator) Apply(key, value string) (string, error) {
	if len(value) <= t.maxBytes {
		return value, nil
	}
	// Walk back from the limit to avoid splitting a multi-byte rune.
	cut := t.maxBytes
	for cut > 0 && !utf8.RuneStart(value[cut]) {
		cut--
	}
	return value[:cut], &ErrTruncated{
		Key:      key,
		Original: len(value),
		Limit:    cut,
	}
}

// Map applies Apply to every entry in env, returning a new map and a slice
// of any truncation errors that occurred.
func (t *Truncator) Map(env map[string]string) (map[string]string, []error) {
	out := make(map[string]string, len(env))
	var errs []error
	for k, v := range env {
		trimmed, err := t.Apply(k, v)
		out[k] = trimmed
		if err != nil {
			errs = append(errs, err)
		}
	}
	return out, errs
}
