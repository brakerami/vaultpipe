// Package evict implements a thread-safe, fixed-capacity least-recently-used
// (LRU) cache intended for bounding the number of secret values held in
// memory by vaultpipe.
//
// When the cache is full, the entry that was least recently read or written is
// silently dropped. An optional Evicted hook lets callers react to removals
// (e.g. to emit metrics or audit log entries).
//
// # Concurrency
//
// All methods on [LRU] are safe for concurrent use by multiple goroutines.
// The Evicted callback, if set, is invoked while the internal lock is NOT
// held, so callers may safely call back into the cache from within the hook
// without deadlocking.
//
// # Usage
//
//	l := evict.New(128)
//	l.Evicted = func(k, v string) { log.Printf("evicted %s", k) }
//	l.Set("secret/db#password", "s3cr3t")
//	v, ok := l.Get("secret/db#password")
package evict
