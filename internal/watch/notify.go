package watch

import (
	"context"
	"fmt"
	"io"
	"time"
)

// Event describes a renewal notification.
type Event struct {
	Ref       string
	Triggered time.Time
	Err       error
}

// String returns a human-readable representation of the event.
func (e Event) String() string {
	if e.Err != nil {
		return fmt.Sprintf("watch: renewal failed for %q: %v", e.Ref, e.Err)
	}
	return fmt.Sprintf("watch: renewed %q at %s", e.Ref, e.Triggered.Format(time.RFC3339))
}

// LoggingRenewFunc returns a RenewFunc that wraps inner and writes
// structured events to w after each renewal attempt.
func LoggingRenewFunc(inner RenewFunc, w io.Writer) RenewFunc {
	return func(ctx context.Context, ref string) error {
		err := inner(ctx, ref)
		ev := Event{Ref: ref, Triggered: time.Now(), Err: err}
		_, _ = fmt.Fprintln(w, ev.String())
		return err
	}
}

// ChannelRenewFunc returns a RenewFunc that wraps inner and sends
// Events on ch (non-blocking). ch must be buffered by the caller.
func ChannelRenewFunc(inner RenewFunc, ch chan<- Event) RenewFunc {
	return func(ctx context.Context, ref string) error {
		err := inner(ctx, ref)
		ev := Event{Ref: ref, Triggered: time.Now(), Err: err}
		select {
		case ch <- ev:
		default:
		}
		return err
	}
}
