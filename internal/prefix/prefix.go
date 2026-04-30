// Package prefix provides a middleware that applies a namespace prefix
// to secret paths before they are resolved from Vault.
package prefix

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// FetchFunc resolves a secret value for the given path.
type FetchFunc func(ctx context.Context, path string) (string, error)

// Prefixer prepends a fixed namespace to every secret path.
type Prefixer struct {
	prefix string
}

// New returns a Prefixer that prepends prefix to all paths.
// prefix must be non-empty and must not contain consecutive slashes.
func New(prefix string) (*Prefixer, error) {
	if prefix == "" {
		return nil, errors.New("prefix: prefix must not be empty")
	}
	if strings.Contains(prefix, "//") {
		return nil, fmt.Errorf("prefix: prefix %q contains consecutive slashes", prefix)
	}
	return &Prefixer{prefix: strings.TrimRight(prefix, "/")}, nil
}

// Apply returns the fully-qualified path by joining the prefix and path.
// A single slash is always used as the separator.
func (p *Prefixer) Apply(path string) string {
	path = strings.TrimLeft(path, "/")
	if path == "" {
		return p.prefix
	}
	return p.prefix + "/" + path
}

// Wrap returns a FetchFunc that rewrites every path through Apply before
// delegating to next.
func (p *Prefixer) Wrap(next FetchFunc) FetchFunc {
	return func(ctx context.Context, path string) (string, error) {
		return next(ctx, p.Apply(path))
	}
}

// String returns the configured prefix.
func (p *Prefixer) String() string { return p.prefix }
