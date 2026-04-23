// Package window provides a sliding-window counter and rate tracker.
//
// Window tracks discrete events over a rolling time interval divided into
// fixed-size buckets. Stale buckets are discarded automatically when the
// window advances, giving an approximate but memory-efficient count.
//
// Rate wraps Window to expose per-second event rate calculations and a
// simple threshold check useful for adaptive back-pressure or alerting.
//
// Example:
//
//	r := window.NewRate(10*time.Second, 20)
//	r.Observe(1)
//	if r.Exceeds(100) {
//		// shed load
//	}
package window
