// Package rotate detects secret rotation by comparing newly resolved Vault
// secret values against previously observed ones.
//
// Usage:
//
//	det := rotate.New()
//	det.Seed("DB_PASSWORD", initialValue)
//	det.OnChange(func(ctx context.Context, key, old, new string) {
//		log.Printf("secret %s rotated", key)
//	})
//
//	// Later, after re-fetching:
//	det.Observe(ctx, "DB_PASSWORD", freshValue)
//
// Hooks are invoked synchronously in the calling goroutine. If concurrent
// safety is required callers should ensure Observe is not called from
// multiple goroutines simultaneously, or accept that hooks may interleave.
package rotate
