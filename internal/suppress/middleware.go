package suppress

import (
	"context"
	"fmt"
)

// FetchFunc is a function that retrieves a secret value by path.
type FetchFunc func(ctx context.Context, path string) (string, error)

// GuardedFetch wraps a FetchFunc so that fetch errors for the same
// path are only forwarded to the caller once per suppression window.
// Subsequent errors within the window return the last error silently
// by returning an empty string and nil — callers should treat an empty
// result as a suppressed failure and retain the previously cached value.
func GuardedFetch(s *Suppressor, next FetchFunc) FetchFunc {
	return func(ctx context.Context, path string) (string, error) {
		val, err := next(ctx, path)
		if err != nil {
			key := fmt.Sprintf("err:%s", path)
			if s.Allow(key) {
				return "", err
			}
			// suppressed — return empty, no error
			return "", nil
		}
		return val, nil
	}
}
