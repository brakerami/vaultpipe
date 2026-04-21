// Package ratelimit implements a token-bucket rate limiter used to throttle
// outbound requests to the Vault API.
//
// When vaultpipe resolves a large number of secret references at startup,
// it can generate a burst of concurrent Vault reads. The Limiter ensures
// requests are spread over time, respecting Vault's own rate limits and
// avoiding unnecessary 429 responses.
//
// Usage:
//
//	limiter, err := ratelimit.New(20) // 20 requests per second
//	if err != nil { ... }
//
//	if err := limiter.Wait(ctx); err != nil {
//	    return err // context cancelled
//	}
//	// safe to call Vault now
package ratelimit
