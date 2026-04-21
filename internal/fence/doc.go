// Package fence implements a per-key mutual-exclusion barrier for
// concurrent secret resolution.
//
// When multiple goroutines attempt to resolve the same Vault path at
// the same time, only the first acquires the fence token and performs
// the actual fetch. All other callers block in Wait until the token is
// released, then proceed to call the underlying resolver themselves
// (relying on the upstream cache to serve the now-warm result cheaply).
//
// Usage:
//
//	f := fence.New()
//	wrapped := fence.Deduplicate(f, reg, myResolveFunc)
//	value, err := wrapped(ctx, "secret/data/myapp#password")
//
// The Fence is safe for concurrent use by multiple goroutines.
package fence
