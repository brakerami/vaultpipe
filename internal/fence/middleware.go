package fence

import (
	"context"
	"fmt"

	"github.com/yourusername/vaultpipe/internal/metrics"
)

// ResolveFunc is any function that fetches a secret value by path.
type ResolveFunc func(ctx context.Context, path string) (string, error)

// Deduplicate wraps a ResolveFunc so that concurrent calls for the
// same path are collapsed: the first caller performs the fetch while
// subsequent callers wait for it to finish, then retry via next.
// This prevents thundering-herd spikes against Vault.
func Deduplicate(f *Fence, reg *metrics.Registry, next ResolveFunc) ResolveFunc {
	return func(ctx context.Context, path string) (string, error) {
		release, err := f.Acquire(ctx, path)
		if err == ErrFenced {
			// Another goroutine is already fetching; wait then delegate.
			if werr := f.Wait(ctx, path); werr != nil {
				return "", fmt.Errorf("fence wait: %w", werr)
			}
			// Gate is now free — call next directly (result may be cached
			// upstream; fence does not cache itself).
			if reg != nil {
				reg.Counter("fence.deduplicated").Inc()
			}
			return next(ctx, path)
		}
		if err != nil {
			return "", err
		}
		defer release()
		return next(ctx, path)
	}
}
