package batch

import (
	"errors"
	"fmt"
	"strings"
)

// MultiError aggregates one or more errors produced during a batch resolve.
type MultiError struct {
	Errs []error
}

// Error implements the error interface.
func (m *MultiError) Error() string {
	msgs := make([]string, len(m.Errs))
	for i, e := range m.Errs {
		msgs[i] = e.Error()
	}
	return fmt.Sprintf("batch: %d error(s): %s", len(m.Errs), strings.Join(msgs, "; "))
}

// Unwrap returns the slice of wrapped errors for use with errors.As / errors.Is.
func (m *MultiError) Unwrap() []error { return m.Errs }

// Collect converts a slice of Results into a MultiError when any result
// carries a non-nil error. Returns nil when all results succeeded.
func Collect(results []Result) error {
	var errs []error
	for _, r := range results {
		if r.Err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", r.Key, r.Err))
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return &MultiError{Errs: errs}
}

// IsPartialFailure reports whether err is a *MultiError.
func IsPartialFailure(err error) bool {
	var me *MultiError
	return errors.As(err, &me)
}
