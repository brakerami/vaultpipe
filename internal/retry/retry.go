// Package retry provides configurable retry logic with exponential backoff
// for transient errors when fetching secrets from Vault.
package retry

import (
	"context"
	"errors"
	"math"
	"time"
)

// ErrMaxAttempts is returned when all retry attempts are exhausted.
var ErrMaxAttempts = errors.New("retry: max attempts reached")

// Config holds retry policy parameters.
type Config struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
	Multiplier  float64
}

// DefaultConfig returns a sensible default retry configuration.
func DefaultConfig() Config {
	return Config{
		MaxAttempts: 3,
		BaseDelay:   200 * time.Millisecond,
		MaxDelay:    5 * time.Second,
		Multiplier:  2.0,
	}
}

// Do executes fn up to cfg.MaxAttempts times, backing off between attempts.
// It stops early if ctx is cancelled or fn returns a non-retryable error.
func Do(ctx context.Context, cfg Config, fn func() error) error {
	if cfg.MaxAttempts <= 0 {
		cfg.MaxAttempts = 1
	}
	if cfg.Multiplier <= 1 {
		cfg.Multiplier = 2.0
	}

	var lastErr error
	for attempt := 0; attempt < cfg.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		lastErr = fn()
		if lastErr == nil {
			return nil
		}
		var nr *NonRetryableError
		if errors.As(lastErr, &nr) {
			return nr.Unwrap()
		}
		if attempt < cfg.MaxAttempts-1 {
			delay := backoff(cfg.BaseDelay, cfg.MaxDelay, cfg.Multiplier, attempt)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
		}
	}
	return errors.Join(ErrMaxAttempts, lastErr)
}

// backoff calculates the delay for a given attempt using exponential backoff.
func backoff(base, max time.Duration, multiplier float64, attempt int) time.Duration {
	d := float64(base) * math.Pow(multiplier, float64(attempt))
	if d > float64(max) {
		d = float64(max)
	}
	return time.Duration(d)
}

// NonRetryableError wraps an error to signal that retry should stop immediately.
type NonRetryableError struct {
	cause error
}

// Permanent wraps err so that Do will not retry it.
func Permanent(err error) error {
	if err == nil {
		return nil
	}
	return &NonRetryableError{cause: err}
}

func (e *NonRetryableError) Error() string { return e.cause.Error() }
func (e *NonRetryableError) Unwrap() error { return e.cause }
