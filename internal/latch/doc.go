// Package latch provides a one-shot boolean gate.
//
// A Latch starts in the open (unset) state and transitions to the closed
// (set) state exactly once. After being set it remains set permanently.
//
// Typical use cases inside vaultpipe:
//
//   - Signal that initial secret resolution is complete before allowing the
//     child process to start.
//   - Coordinate between a background watcher goroutine and the main pipeline
//     so that the pipeline waits until the first successful fetch.
//
// Example:
//
//	ready := latch.New()
//	go func() {
//		// fetch secrets …
//		ready.Set()
//	}()
//	if err := ready.Wait(ctx); err != nil {
//		return err
//	}
package latch
