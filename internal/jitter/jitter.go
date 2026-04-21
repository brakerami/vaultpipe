// Package jitter provides utilities for adding randomised delay to
// retry and renewal schedules, reducing thundering-herd effects when
// many secrets share the same TTL.
package jitter

import (
	"math/rand"
	"time"
)

// Source is the interface used to obtain random floats, allowing tests
// to inject a deterministic generator.
type Source interface {
	Float64() float64
}

// defaultSource wraps the package-level rand functions.
type defaultSource struct{}

func (defaultSource) Float64() float64 { return rand.Float64() }

// Jitter adds a random fraction of base to base, returning a duration
// in the range [base, base*(1+factor)]. factor is clamped to [0, 1].
func Jitter(base time.Duration, factor float64) time.Duration {
	return JitterWith(base, factor, defaultSource{})
}

// JitterWith is the same as Jitter but accepts an explicit Source so
// callers can supply a seeded or deterministic generator.
func JitterWith(base time.Duration, factor float64, src Source) time.Duration {
	if factor < 0 {
		factor = 0
	}
	if factor > 1 {
		factor = 1
	}
	if base <= 0 {
		return base
	}
	delta := float64(base) * factor * src.Float64()
	return base + time.Duration(delta)
}

// Full returns a duration chosen uniformly at random from [0, base].
// This is the "full jitter" strategy recommended for back-off loops.
func Full(base time.Duration) time.Duration {
	return FullWith(base, defaultSource{})
}

// FullWith is Full with an explicit Source.
func FullWith(base time.Duration, src Source) time.Duration {
	if base <= 0 {
		return 0
	}
	return time.Duration(float64(base) * src.Float64())
}

// Equal returns a duration chosen uniformly from [base/2, base].
// This is the "equal jitter" strategy that guarantees at least half
// of the base delay is always observed.
func Equal(base time.Duration) time.Duration {
	return EqualWith(base, defaultSource{})
}

// EqualWith is Equal with an explicit Source.
func EqualWith(base time.Duration, src Source) time.Duration {
	if base <= 0 {
		return 0
	}
	half := base / 2
	return half + FullWith(half, src)
}
