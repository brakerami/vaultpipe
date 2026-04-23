// Package prefetch schedules background secret pre-loading so that cached
// values are refreshed before their leases expire.
//
// Usage:
//
//	p := prefetch.New(vaultFetch, func(path, value string) {
//		cache.Set(path, value)
//	})
//
//	// Refresh when 20 % of the 5-minute TTL remains (i.e. after 4 minutes).
//	p.Register("kv/data/db-password", 5*time.Minute, 0.2)
//
//	// Cancel a single entry.
//	p.Deregister("kv/data/db-password")
//
//	// Cancel everything on shutdown.
//	p.Stop()
package prefetch
