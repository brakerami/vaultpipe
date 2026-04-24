package shadow

import (
	"context"
	"fmt"
)

// FetchFunc fetches the current value for a secret path.
type FetchFunc func(ctx context.Context, path string) (string, error)

// TrackingFetch wraps next so that every successfully fetched value is
// recorded in store. If drift is detected (the new value differs from the
// previously shadowed value), onDrift is called with the affected key.
func TrackingFetch(store *Shadow, onDrift func(key, old, new string), next FetchFunc) FetchFunc {
	return func(ctx context.Context, path string) (string, error) {
		value, err := next(ctx, path)
		if err != nil {
			return "", err
		}

		if store.Drifted(path, value) {
			if onDrift != nil {
				old, _ := store.Get(path)
				onDrift(path, old.Value, value)
			}
		}

		store.Set(path, value)
		return value, nil
	}
}

// ErrDriftDetected is returned by StrictFetch when drift is detected.
type ErrDriftDetected struct {
	Key string
	Old string
	New string
}

func (e *ErrDriftDetected) Error() string {
	return fmt.Sprintf("shadow: drift detected for %q: old=%q new=%q", e.Key, e.Old, e.New)
}

// StrictFetch wraps next and returns ErrDriftDetected instead of the new value
// when drift is observed, leaving the shadow unchanged.
func StrictFetch(store *Shadow, next FetchFunc) FetchFunc {
	return func(ctx context.Context, path string) (string, error) {
		value, err := next(ctx, path)
		if err != nil {
			return "", err
		}

		if store.Drifted(path, value) {
			old, _ := store.Get(path)
			return "", &ErrDriftDetected{Key: path, Old: old.Value, New: value}
		}

		store.Set(path, value)
		return value, nil
	}
}
