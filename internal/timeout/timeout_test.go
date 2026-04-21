package timeout_test

import (
	"context"
	"testing"
	"time"

	"github.com/your-org/vaultpipe/internal/timeout"
)

func TestDefault_Values(t *testing.T) {
	cfg := timeout.Default()
	if cfg.Fetch != timeout.DefaultFetch {
		t.Fatalf("expected Fetch=%v, got %v", timeout.DefaultFetch, cfg.Fetch)
	}
	if cfg.Renew != timeout.DefaultRenew {
		t.Fatalf("expected Renew=%v, got %v", timeout.DefaultRenew, cfg.Renew)
	}
}

func TestValidate_Valid(t *testing.T) {
	if err := timeout.Default().Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_TooShort(t *testing.T) {
	cfg := timeout.Config{Fetch: 10 * time.Millisecond, Renew: timeout.DefaultRenew}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for too-short fetch, got nil")
	}
}

func TestValidate_TooLong(t *testing.T) {
	cfg := timeout.Config{Fetch: timeout.DefaultFetch, Renew: 10 * time.Minute}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for too-long renew, got nil")
	}
}

func TestWithFetch_CancelsAfterDeadline(t *testing.T) {
	cfg := timeout.Config{Fetch: 50 * time.Millisecond, Renew: timeout.DefaultRenew}
	ctx, cancel := cfg.WithFetch(context.Background())
	defer cancel()

	select {
	case <-ctx.Done():
		// expected
	case <-time.After(500 * time.Millisecond):
		t.Fatal("context was not cancelled within expected window")
	}
}

func TestWithRenew_CancelsAfterDeadline(t *testing.T) {
	cfg := timeout.Config{Fetch: timeout.DefaultFetch, Renew: 50 * time.Millisecond}
	ctx, cancel := cfg.WithRenew(context.Background())
	defer cancel()

	select {
	case <-ctx.Done():
		// expected
	case <-time.After(500 * time.Millisecond):
		t.Fatal("context was not cancelled within expected window")
	}
}

func TestIsTimeout_DeadlineExceeded(t *testing.T) {
	if !timeout.IsTimeout(context.DeadlineExceeded) {
		t.Fatal("expected IsTimeout=true for context.DeadlineExceeded")
	}
}

func TestIsTimeout_ErrTimeout(t *testing.T) {
	if !timeout.IsTimeout(timeout.ErrTimeout) {
		t.Fatal("expected IsTimeout=true for ErrTimeout")
	}
}

func TestIsTimeout_OtherError(t *testing.T) {
	if timeout.IsTimeout(context.Canceled) {
		t.Fatal("expected IsTimeout=false for context.Canceled")
	}
}
