// Package renew provides automatic secret renewal for Vault leases.
// It coordinates between the watch package (which tracks TTLs) and the
// resolver package (which fetches secrets) to keep secrets fresh.
package renew

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/yourusername/vaultpipe/internal/cache"
	"github.com/yourusername/vaultpipe/internal/resolver"
	"github.com/yourusername/vaultpipe/internal/watch"
)

// RenewFunc is called when a secret has been renewed with the new value.
type RenewFunc func(key, value string)

// Manager coordinates secret renewal using a Watcher and Resolver.
type Manager struct {
	watcher  *watch.Watcher
	resolver *resolver.Resolver
	cache    *cache.Cache
	logger   *slog.Logger
	onRenew  RenewFunc

	mu      sync.Mutex
	pathMap map[string]string // maps watch key -> vault path
}

// Config holds options for creating a Manager.
type Config struct {
	Watcher  *watch.Watcher
	Resolver *resolver.Resolver
	Cache    *cache.Cache
	Logger   *slog.Logger
	OnRenew  RenewFunc
}

// New creates a new renewal Manager.
func New(cfg Config) (*Manager, error) {
	if cfg.Watcher == nil {
		return nil, fmt.Errorf("renew: watcher is required")
	}
	if cfg.Resolver == nil {
		return nil, fmt.Errorf("renew: resolver is required")
	}

	logger := cfg.Logger
	if logger == nil {
		logger = slog.Default()
	}

	return &Manager{
		watcher:  cfg.Watcher,
		resolver: cfg.Resolver,
		cache:    cfg.Cache,
		logger:   logger,
		onRenew:  cfg.OnRenew,
		pathMap:  make(map[string]string),
	}, nil
}

// Track registers a secret path for renewal tracking with the given TTL.
// The key is an identifier used to correlate renewal events back to the
// environment variable or config entry that owns the secret.
func (m *Manager) Track(ctx context.Context, key, vaultPath string, ttl time.Duration) {
	m.mu.Lock()
	m.pathMap[key] = vaultPath
	m.mu.Unlock()

	renewFn := watch.LoggingRenewFunc(m.logger, key, func(ctx context.Context) error {
		return m.renewSecret(ctx, key)
	})

	m.watcher.Add(key, ttl, renewFn)
	m.logger.Info("tracking secret for renewal",
		"key", key,
		"path", vaultPath,
		"ttl", ttl,
	)
}

// Untrack removes a secret from renewal tracking.
func (m *Manager) Untrack(key string) {
	m.mu.Lock()
	delete(m.pathMap, key)
	m.mu.Unlock()
	m.watcher.Remove(key)
}

// TrackedKeys returns a snapshot of all currently tracked secret keys.
func (m *Manager) TrackedKeys() []string {
	m.mu.Lock()
	defer m.mu.Unlock()

	keys := make([]string, 0, len(m.pathMap))
	for k := range m.pathMap {
		keys = append(keys, k)
	}
	return keys
}

// renewSecret fetches a fresh value for the secret identified by key.
func (m *Manager) renewSecret(ctx context.Context, key string) error {
	m.mu.Lock()
	path, ok := m.pathMap[key]
	m.mu.Unlock()

	if !ok {
		return fmt.Errorf("renew: unknown key %q", key)
	}

	value, err := m.resolver.Resolve(ctx, path)
	if err != nil {
		return fmt.Errorf("renew: failed to resolve %q: %w", path, err)
	}

	if m.cache != nil {
		m.cache.Invalidate(path)
	}

	if m.onRenew != nil {
		m.onRenew(key, value)
	}

	return nil
}

// Stop shuts down the renewal manager and its watcher.
func (m *Manager) Stop() {
	m.watcher.Stop()
}
