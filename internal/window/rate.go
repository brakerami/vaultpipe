package window

import (
	"fmt"
	"time"
)

// Rate wraps a Window and exposes per-second event rate helpers.
type Rate struct {
	w *Window
}

// NewRate returns a Rate backed by a Window of the given size and bucket count.
func NewRate(size time.Duration, buckets int) *Rate {
	return &Rate{w: New(size, buckets)}
}

// Observe records n events.
func (r *Rate) Observe(n int64) {
	r.w.Add(n)
}

// PerSecond returns the average number of events per second within the window.
func (r *Rate) PerSecond() float64 {
	count := r.w.Count()
	seconds := r.w.size.Seconds()
	if seconds == 0 {
		return 0
	}
	return float64(count) / seconds
}

// Exceeds reports whether the current per-second rate is above the given limit.
func (r *Rate) Exceeds(limit float64) bool {
	return r.PerSecond() > limit
}

// String returns a human-readable summary.
func (r *Rate) String() string {
	return fmt.Sprintf("%.2f events/s (window=%s)", r.PerSecond(), r.w.size)
}
