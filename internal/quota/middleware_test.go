package quota_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"vaultpipe/internal/quota"
)

func TestGuardedFetch_AllowsUnderQuota(t *testing.T) {
	q, _ := quota.New(quota.Config{MaxRequests: 5, Window: time.Minute})

	next := func(_ context.Context, path string) (string, error) {
		return "s3cr3t", nil
	}

	fetch := quota.GuardedFetch(q, next)
	val, err := fetch(context.Background(), "secret/mykey")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "s3cr3t" {
		t.Fatalf("expected s3cr3t, got %q", val)
	}
}

func TestGuardedFetch_BlocksOverQuota(t *testing.T) {
	q, _ := quota.New(quota.Config{MaxRequests: 2, Window: time.Minute})

	called := 0
	next := func(_ context.Context, path string) (string, error) {
		called++
		return "value", nil
	}

	fetch := quota.GuardedFetch(q, next)
	_, _ = fetch(context.Background(), "secret/limited")
	_, _ = fetch(context.Background(), "secret/limited")
	_, err := fetch(context.Background(), "secret/limited")

	if !errors.Is(err, quota.ErrExceeded) {
		t.Fatalf("expected ErrExceeded, got %v", err)
	}
	if called != 2 {
		t.Fatalf("expected next called 2 times, got %d", called)
	}
}

func TestGuardedFetch_PropagatesNextError(t *testing.T) {
	q, _ := quota.New(quota.Config{MaxRequests: 5, Window: time.Minute})

	expected := errors.New("vault unavailable")
	next := func(_ context.Context, _ string) (string, error) {
		return "", expected
	}

	fetch := quota.GuardedFetch(q, next)
	_, err := fetch(context.Background(), "secret/broken")
	if !errors.Is(err, expected) {
		t.Fatalf("expected underlying error, got %v", err)
	}
}
