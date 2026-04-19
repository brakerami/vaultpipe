// Package cache provides an in-memory TTL cache for resolved secrets
// to avoid redundant Vault API calls within a single vaultpipe run.
package cache

import (
	"sync"
	"time"
)

// Entry holds a cached secret value and its expiry.
type Entry struct {
	Value     string
	FetchedAt time.Time
	ExpiresAt time.Time
}

// Expired returns true if the entry is past its TTL.
func (e Entry) Expired() bool {
	return time.Now().After(e.ExpiresAt)
}

// Cache is a thread-safe in-memory store for secret values.
type Cache struct {
	mu      sync.RWMutex
	entries map[string]Entry
	ttl     time.Duration
}

// New creates a Cache with the given TTL. A zero TTL disables caching.
func New(ttl time.Duration) *Cache {
	return &Cache{
		entries: make(map[string]Entry),
		ttl:     ttl,
	}
}

// Get retrieves a cached value by key. Returns the value and whether it was
// found and still valid.
func (c *Cache) Get(key string) (string, bool) {
	if c.ttl == 0 {
		return "", false
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.entries[key]
	if !ok || e.Expired() {
		return "", false
	}
	return e.Value, true
}

// Set stores a value under key with the configured TTL.
func (c *Cache) Set(key, value string) {
	if c.ttl == 0 {
		return
	}
	now := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = Entry{
		Value:     value,
		FetchedAt: now,
		ExpiresAt: now.Add(c.ttl),
	}
}

// Invalidate removes a single key from the cache.
func (c *Cache) Invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}

// Len returns the number of entries currently in the cache (including expired).
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}
