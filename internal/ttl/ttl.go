// Package ttl provides utilities for parsing and validating
// time-to-live durations used in secret caching and lease management.
package ttl

import (
	"errors"
	"fmt"
	"time"
)

// Default durations used when no explicit TTL is specified.
const (
	DefaultTTL = 5 * time.Minute
	MinTTL     = 1 * time.Second
	MaxTTL     = 24 * time.Hour
)

// ErrTTLTooShort is returned when the parsed duration is below MinTTL.
var ErrTTLTooShort = errors.New("ttl: duration is below minimum allowed value")

// ErrTTLTooLong is returned when the parsed duration exceeds MaxTTL.
var ErrTTLTooLong = errors.New("ttl: duration exceeds maximum allowed value")

// Parse parses a duration string and validates it falls within the
// acceptable range [MinTTL, MaxTTL]. An empty string returns DefaultTTL.
func Parse(s string) (time.Duration, error) {
	if s == "" {
		return DefaultTTL, nil
	}

	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, fmt.Errorf("ttl: invalid duration %q: %w", s, err)
	}

	return Validate(d)
}

// Validate checks that d is within [MinTTL, MaxTTL].
func Validate(d time.Duration) (time.Duration, error) {
	if d < MinTTL {
		return 0, ErrTTLTooShort
	}
	if d > MaxTTL {
		return 0, ErrTTLTooLong
	}
	return d, nil
}

// MustParse is like Parse but panics on error. Intended for tests and
// compile-time constants only.
func MustParse(s string) time.Duration {
	d, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return d
}
