package quota

import (
	"context"
	"fmt"
)

// FetchFunc is the signature of a function that fetches a secret value by path.
type FetchFunc func(ctx context.Context, path string) (string, error)

// GuardedFetch wraps a FetchFunc with quota enforcement.
// If the key exceeds its quota, the fetch is not performed and ErrExceeded is returned.
func GuardedFetch(q *Quota, next FetchFunc) FetchFunc {
	return func(ctx context.Context, path string) (string, error) {
		if err := q.Allow(path); err != nil {
			return "", fmt.Errorf("quota guard: %w", err)
		}
		return next(ctx, path)
	}
}
