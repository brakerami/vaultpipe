package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/your-org/vaultpipe/internal/retry"
)

var errTransient = errors.New("transient error")

func TestDo_SucceedsOnFirstAttempt(t *testing.T) {
	calls := 0
	err := retry.Do(context.Background(), retry.DefaultConfig(), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_RetriesOnTransientError(t *testing.T) {
	calls := 0
	cfg := retry.Config{MaxAttempts: 3, BaseDelay: time.Millisecond, MaxDelay: 10 * time.Millisecond, Multiplier: 2}
	err := retry.Do(context.Background(), cfg, func() error {
		calls++
		if calls < 3 {
			return errTransient
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil after retries, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ExhaustsMaxAttempts(t *testing.T) {
	calls := 0
	cfg := retry.Config{MaxAttempts: 2, BaseDelay: time.Millisecond, MaxDelay: 5 * time.Millisecond, Multiplier: 2}
	err := retry.Do(context.Background(), cfg, func() error {
		calls++
		return errTransient
	})
	if !errors.Is(err, retry.ErrMaxAttempts) {
		t.Fatalf("expected ErrMaxAttempts, got %v", err)
	}
	if calls != 2 {
		t.Fatalf("expected 2 calls, got %d", calls)
	}
}

func TestDo_StopsOnPermanentError(t *testing.T) {
	calls := 0
	permErr := errors.New("permission denied")
	cfg := retry.Config{MaxAttempts: 5, BaseDelay: time.Millisecond, MaxDelay: 5 * time.Millisecond, Multiplier: 2}
	err := retry.Do(context.Background(), cfg, func() error {
		calls++
		return retry.Permanent(permErr)
	})
	if !errors.Is(err, permErr) {
		t.Fatalf("expected permErr, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_RespectsContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	calls := 0
	err := retry.Do(ctx, retry.DefaultConfig(), func() error {
		calls++
		return errTransient
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if calls != 0 {
		t.Fatalf("expected 0 calls, got %d", calls)
	}
}

func TestPermanent_NilPassthrough(t *testing.T) {
	if retry.Permanent(nil) != nil {
		t.Fatal("expected nil for Permanent(nil)")
	}
}
