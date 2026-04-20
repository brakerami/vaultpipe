package retry_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/your-org/vaultpipe/internal/retry"
)

// TestDo_ContextTimeoutDuringBackoff ensures that a context deadline exceeded
// during a backoff sleep is respected promptly.
func TestDo_ContextTimeoutDuringBackoff(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	cfg := retry.Config{
		MaxAttempts: 10,
		BaseDelay:   30 * time.Millisecond,
		MaxDelay:    500 * time.Millisecond,
		Multiplier:  2.0,
	}

	var calls int32
	start := time.Now()
	err := retry.Do(ctx, cfg, func() error {
		atomic.AddInt32(&calls, 1)
		return errors.New("transient")
	})
	elapsed := time.Since(start)

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected DeadlineExceeded, got %v", err)
	}
	if elapsed > 200*time.Millisecond {
		t.Fatalf("took too long: %v", elapsed)
	}
	if atomic.LoadInt32(&calls) == 0 {
		t.Fatal("expected at least one call")
	}
}

// TestDo_SucceedsOnSecondAttempt verifies the happy-path retry scenario.
func TestDo_SucceedsOnSecondAttempt(t *testing.T) {
	cfg := retry.Config{
		MaxAttempts: 3,
		BaseDelay:   time.Millisecond,
		MaxDelay:    10 * time.Millisecond,
		Multiplier:  2.0,
	}

	attempts := 0
	err := retry.Do(context.Background(), cfg, func() error {
		attempts++
		if attempts == 1 {
			return errors.New("first attempt fails")
		}
		return nil
	})

	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if attempts != 2 {
		t.Fatalf("expected 2 attempts, got %d", attempts)
	}
}
