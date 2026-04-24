package sieve

import (
	"context"
	"fmt"
)

// FetchFunc is the signature of a secret-fetching function.
type FetchFunc func(ctx context.Context, path string) (map[string]string, error)

// GuardedFetch wraps next so that any call for a denied path is rejected
// before reaching Vault. Permitted paths are passed through unchanged.
//
//	guarded := sieve.GuardedFetch(s, client.Fetch)
//	values, err := guarded(ctx, "secret/data/myapp")
func GuardedFetch(s *Sieve, next FetchFunc) FetchFunc {
	if s == nil {
		panic("sieve: GuardedFetch requires a non-nil Sieve")
	}
	if next == nil {
		panic("sieve: GuardedFetch requires a non-nil FetchFunc")
	}
	return func(ctx context.Context, path string) (map[string]string, error) {
		if err := s.Check(path); err != nil {
			return nil, fmt.Errorf("sieve: access denied for %q: %w", path, err)
		}
		return next(ctx, path)
	}
}
