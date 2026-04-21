// Package debounce coalesces rapid successive function calls into a single
// execution after a configurable quiet period.
//
// # Overview
//
// When secrets near expiry, multiple watchers may independently observe the
// same lease and schedule a renewal within milliseconds of each other.
// Debouncing ensures that only one renewal reaches Vault per key within the
// configured window, reducing unnecessary API load.
//
// # Usage
//
//	d := debounce.New(500*time.Millisecond, func(ctx context.Context) {
//		// called at most once per 500 ms quiet window
//	})
//	d.Trigger(ctx) // resets the timer on every call
//	d.Stop()       // cancel without executing
//	d.Flush(ctx)   // execute immediately and cancel timer
package debounce
