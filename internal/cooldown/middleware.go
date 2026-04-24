package cooldown

import (
	"context"
	"fmt"
	"time"
)

// FetchFunc is the signature for a secret-fetching operation.
type FetchFunc func(ctx context.Context, path string) (string, error)

// ErrCooldown is returned by GuardedFetch when the cooldown period for
// the requested path has not yet elapsed.
type ErrCooldown struct {
	Path      string
	Remaining time.Duration
}

func (e *ErrCooldown) Error() string {
	return fmt.Sprintf("cooldown: %q is rate-limited for another %v", e.Path, e.Remaining)
}

// GuardedFetch wraps next with a per-path cooldown. If the cooldown for
// path has not expired, GuardedFetch returns an *ErrCooldown without
// calling next. Otherwise it delegates to next and, on success, records
// the attempt so future calls within the interval are suppressed.
func GuardedFetch(c *Cooldown, next FetchFunc) FetchFunc {
	return func(ctx context.Context, path string) (string, error) {
		if !c.Allow(path) {
			return "", &ErrCooldown{
				Path:      path,
				Remaining: c.Remaining(path),
			}
		}
		val, err := next(ctx, path)
		if err != nil {
			// Roll back so a transient error does not lock out retries.
			c.Reset(path)
			return "", err
		}
		return val, nil
	}
}
