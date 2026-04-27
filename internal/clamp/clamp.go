// Package clamp provides value clamping utilities for secret metadata fields
// such as TTLs, retry counts, and concurrency limits.
package clamp

import (
	"fmt"
	"time"
)

// Int clamps v to the inclusive range [min, max].
// It returns an error if min > max.
func Int(v, min, max int) (int, error) {
	if min > max {
		return 0, fmt.Errorf("clamp: min %d exceeds max %d", min, max)
	}
	if v < min {
		return min, nil
	}
	if v > max {
		return max, nil
	}
	return v, nil
}

// MustInt is like Int but panics on invalid bounds.
func MustInt(v, min, max int) int {
	out, err := Int(v, min, max)
	if err != nil {
		panic(err)
	}
	return out
}

// Duration clamps d to the inclusive range [min, max].
// It returns an error if min > max.
func Duration(d, min, max time.Duration) (time.Duration, error) {
	if min > max {
		return 0, fmt.Errorf("clamp: min %s exceeds max %s", min, max)
	}
	if d < min {
		return min, nil
	}
	if d > max {
		return max, nil
	}
	return d, nil
}

// MustDuration is like Duration but panics on invalid bounds.
func MustDuration(d, min, max time.Duration) time.Duration {
	out, err := Duration(d, min, max)
	if err != nil {
		panic(err)
	}
	return out
}
