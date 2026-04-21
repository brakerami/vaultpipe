// Package backoff provides exponential backoff with jitter for retry logic.
package backoff

import (
	"math"
	"math/rand"
	"time"
)

// Config holds parameters for exponential backoff calculation.
type Config struct {
	// Base is the initial delay before the first retry.
	Base time.Duration
	// Max caps the computed delay.
	Max time.Duration
	// Multiplier is applied to the previous delay on each attempt.
	Multiplier float64
	// Jitter adds random noise as a fraction of the computed delay (0–1).
	Jitter float64
}

// Default returns a Config suitable for most Vault API interactions.
func Default() Config {
	return Config{
		Base:       200 * time.Millisecond,
		Max:        30 * time.Second,
		Multiplier: 2.0,
		Jitter:     0.2,
	}
}

// Next computes the delay for the given attempt number (zero-indexed).
// It applies exponential growth capped at Max and adds optional jitter.
func (c Config) Next(attempt int) time.Duration {
	if c.Multiplier < 1.0 {
		c.Multiplier = 2.0
	}

	exp := math.Pow(c.Multiplier, float64(attempt))
	delay := time.Duration(float64(c.Base) * exp)

	if delay > c.Max {
		delay = c.Max
	}

	if c.Jitter > 0 {
		// nolint:gosec — jitter does not require cryptographic randomness
		jitter := time.Duration(float64(delay) * c.Jitter * rand.Float64())
		delay += jitter
	}

	return delay
}
