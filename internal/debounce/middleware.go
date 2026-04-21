package debounce

import (
	"context"
	"time"
)

// RenewFunc is a function that renews a secret lease identified by key.
type RenewFunc func(ctx context.Context, key string) error

// DebouncedRenewer wraps a RenewFunc so that renewal calls for the same key
// are coalesced within the debounce window. This prevents thundering-herd
// behaviour when multiple watchers observe the same lease nearing expiry.
type DebouncedRenewer struct {
	wait     time.Duration
	renew    RenewFunc
	debouncers map[string]*Debouncer
}

// NewDebouncedRenewer returns a DebouncedRenewer that delays renewal of each
// key by wait, coalescing duplicate triggers within that window.
func NewDebouncedRenewer(wait time.Duration, fn RenewFunc) *DebouncedRenewer {
	if fn == nil {
		panic("debounce: RenewFunc must not be nil")
	}
	return &DebouncedRenewer{
		wait:       wait,
		renew:      fn,
		debouncers: make(map[string]*Debouncer),
	}
}

// Trigger schedules a renewal for key. Rapid successive calls within wait
// are collapsed into a single renewal attempt.
func (dr *DebouncedRenewer) Trigger(ctx context.Context, key string) {
	d, ok := dr.debouncers[key]
	if !ok {
		d = New(dr.wait, func(ctx context.Context) {
			_ = dr.renew(ctx, key) //nolint:errcheck // caller observes via audit log
		})
		dr.debouncers[key] = d
	}
	d.Trigger(ctx)
}

// Stop cancels all pending renewals.
func (dr *DebouncedRenewer) Stop() {
	for _, d := range dr.debouncers {
		d.Stop()
	}
}
