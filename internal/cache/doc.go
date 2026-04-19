// Package cache implements a thread-safe, TTL-based in-memory cache used by
// vaultpipe to deduplicate Vault secret lookups within a single process
// invocation. When multiple environment variable mappings reference the same
// Vault path, only the first read incurs a network round-trip; subsequent
// lookups are served from the cache until the TTL elapses.
//
// Caching is entirely optional: passing a zero duration to New disables the
// cache transparently so callers need not special-case the behaviour.
package cache
