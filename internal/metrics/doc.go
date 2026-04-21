// Package metrics provides a lightweight, goroutine-safe counter registry
// for instrumenting vaultpipe internals.
//
// Counters are created on first access and are safe for concurrent use.
// Use [Registry.Snapshot] to obtain a point-in-time view of all counters,
// suitable for logging or exposing via a status endpoint.
//
// Well-known counter names used across the project:
//
//	"secret.fetch.ok"    – successful Vault secret retrievals
//	"secret.fetch.error" – failed Vault secret retrievals
//	"cache.hit"          – secrets served from the in-memory cache
//	"cache.miss"         – secrets not found in cache (Vault was queried)
//	"renew.ok"           – successful lease renewals
//	"renew.error"        – failed lease renewals
package metrics
