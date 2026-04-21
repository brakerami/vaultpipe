// Package timeout provides context-based deadline enforcement for
// secret fetch and renewal operations, with configurable per-operation
// limits and a shared default.
package timeout

import (
	"context"
	"errors"
	"fmt"
	"time"
)

const (
	// DefaultFetch is the default deadline applied to a single secret fetch.
	DefaultFetch = 10 * time.Second
	// DefaultRenew is the default deadline applied to a lease renewal.
	DefaultRenew = 5 * time.Second
	// MinTimeout is the smallest accepted timeout value.
	MinTimeout = 100 * time.Millisecond
	// MaxTimeout is the largest accepted timeout value.
	MaxTimeout = 5 * time.Minute
)

// ErrTimeout is returned when an operation exceeds its deadline.
var ErrTimeout = errors.New("operation timed out")

// Config holds per-operation timeout settings.
type Config struct {
	Fetch time.Duration
	Renew time.Duration
}

// Default returns a Config populated with package-level defaults.
func Default() Config {
	return Config{
		Fetch: DefaultFetch,
		Renew: DefaultRenew,
	}
}

// Validate returns an error if any configured duration is outside the
// accepted [MinTimeout, MaxTimeout] range.
func (c Config) Validate() error {
	for name, d := range map[string]time.Duration{
		"fetch": c.Fetch,
		"renew": c.Renew,
	} {
		if d < MinTimeout || d > MaxTimeout {
			return fmt.Errorf("timeout: %s duration %v out of range [%v, %v]",
				name, d, MinTimeout, MaxTimeout)
		}
	}
	return nil
}

// WithFetch derives a child context that will be cancelled after the
// configured fetch deadline. The caller must invoke the returned
// CancelFunc when the operation completes.
func (c Config) WithFetch(parent context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, c.Fetch)
}

// WithRenew derives a child context that will be cancelled after the
// configured renew deadline.
func (c Config) WithRenew(parent context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, c.Renew)
}

// IsTimeout reports whether err was caused by a deadline exceedance,
// covering both context.DeadlineExceeded and ErrTimeout.
func IsTimeout(err error) bool {
	return errors.Is(err, context.DeadlineExceeded) || errors.Is(err, ErrTimeout)
}
