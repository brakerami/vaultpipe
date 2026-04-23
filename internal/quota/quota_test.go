package quota_test

import (
	"errors"
	"testing"
	"time"

	"vaultpipe/internal/quota"
)

func TestNew_InvalidMaxRequests(t *testing.T) {
	_, err := quota.New(quota.Config{MaxRequests: 0, Window: time.Second})
	if err == nil {
		t.Fatal("expected error for MaxRequests=0")
	}
}

func TestNew_InvalidWindow(t *testing.T) {
	_, err := quota.New(quota.Config{MaxRequests: 1, Window: 0})
	if err == nil {
		t.Fatal("expected error for Window=0")
	}
}

func TestNew_Valid(t *testing.T) {
	q, err := quota.New(quota.Config{MaxRequests: 5, Window: time.Minute})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if q == nil {
		t.Fatal("expected non-nil Quota")
	}
}

func TestAllow_WithinQuota(t *testing.T) {
	q, _ := quota.New(quota.Config{MaxRequests: 3, Window: time.Minute})
	for i := 0; i < 3; i++ {
		if err := q.Allow("secret/foo"); err != nil {
			t.Fatalf("request %d: unexpected error: %v", i+1, err)
		}
	}
}

func TestAllow_ExceedsQuota(t *testing.T) {
	q, _ := quota.New(quota.Config{MaxRequests: 2, Window: time.Minute})
	_ = q.Allow("secret/bar")
	_ = q.Allow("secret/bar")
	err := q.Allow("secret/bar")
	if !errors.Is(err, quota.ErrExceeded) {
		t.Fatalf("expected ErrExceeded, got %v", err)
	}
}

func TestAllow_IndependentKeys(t *testing.T) {
	q, _ := quota.New(quota.Config{MaxRequests: 1, Window: time.Minute})
	if err := q.Allow("key/a"); err != nil {
		t.Fatalf("key/a: unexpected error: %v", err)
	}
	if err := q.Allow("key/b"); err != nil {
		t.Fatalf("key/b: unexpected error: %v", err)
	}
}

func TestReset_ClearsHistory(t *testing.T) {
	q, _ := quota.New(quota.Config{MaxRequests: 1, Window: time.Minute})
	_ = q.Allow("secret/baz")
	q.Reset("secret/baz")
	if err := q.Allow("secret/baz"); err != nil {
		t.Fatalf("after reset: unexpected error: %v", err)
	}
}

func TestCount_ReflectsRequests(t *testing.T) {
	q, _ := quota.New(quota.Config{MaxRequests: 10, Window: time.Minute})
	_ = q.Allow("secret/qux")
	_ = q.Allow("secret/qux")
	if got := q.Count("secret/qux"); got != 2 {
		t.Fatalf("expected count 2, got %d", got)
	}
}

func TestCount_MissingKey_ReturnsZero(t *testing.T) {
	q, _ := quota.New(quota.Config{MaxRequests: 5, Window: time.Minute})
	if got := q.Count("nonexistent"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestAllow_WindowExpiry(t *testing.T) {
	q, _ := quota.New(quota.Config{MaxRequests: 1, Window: 50 * time.Millisecond})
	_ = q.Allow("secret/ttl")
	time.Sleep(60 * time.Millisecond)
	if err := q.Allow("secret/ttl"); err != nil {
		t.Fatalf("after window expiry: unexpected error: %v", err)
	}
}
