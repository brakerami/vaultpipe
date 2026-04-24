// Package batch resolves multiple Vault secret references concurrently.
//
// It provides a Resolver that fans out fetch calls across a bounded worker
// pool, collects results (including partial failures), and surfaces them
// through a MultiError so callers can decide whether to abort or continue
// with the successfully resolved subset.
//
// Typical usage:
//
//	r := batch.New(vaultClient.Read, 8)
//	results := r.Resolve(ctx, refs)
//	if err := batch.Collect(results); err != nil {
//		log.Printf("some secrets failed: %v", err)
//	}
package batch
