// Package backoff implements exponential backoff with optional jitter.
//
// It is intentionally decoupled from retry logic so that callers can
// compute delays independently — for example when scheduling lease
// renewals or rate-limit recovery in the Vault client.
//
// Basic usage:
//
//	cfg := backoff.Default()
//	for attempt := 0; attempt < maxAttempts; attempt++ {
//	    if err := doSomething(); err == nil {
//	        break
//	    }
//	    time.Sleep(cfg.Next(attempt))
//	}
package backoff
