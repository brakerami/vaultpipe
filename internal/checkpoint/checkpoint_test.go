package checkpoint_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/checkpoint"
)

// stubFetch is a controllable fetch function for tests.
type stubFetch struct {
	mu      sync.Mutex
	calls   int
	results []fetchResult
}

type fetchResult struct {
	value string
	err   error
}

func (s *stubFetch) add(value string, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.results = append(s.results, fetchResult{value: value, err: err})
}

func (s *stubFetch) fetch(_ context.Context, key string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.calls++
	if len(s.results) == 0 {
		return "", errors.New("no result configured")
	}
	r := s.results[0]
	if len(s.results) > 1 {
		s.results = s.results[1:]
	}
	return r.value, r.err
}

func TestNew_ReturnsCheckpoint(t *testing.T) {
	cp := checkpoint.New(5, 30*time.Second)
	if cp == nil {
		t.Fatal("expected non-nil checkpoint")
	}
}

func TestLoad_FetchesAndCaches(t *testing.T) {
	stub := &stubFetch{}
	stub.add("secret-value", nil)

	cp := checkpoint.New(5, 30*time.Second)
	val, err := cp.Load(context.Background(), "mykey", stub.fetch)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "secret-value" {
		t.Errorf("got %q, want %q", val, "secret-value")
	}
	if stub.calls != 1 {
		t.Errorf("expected 1 fetch call, got %d", stub.calls)
	}

	// Second load should use cache — no additional fetch.
	stub.add("other-value", nil)
	val2, err := cp.Load(context.Background(), "mykey", stub.fetch)
	if err != nil {
		t.Fatalf("unexpected error on second load: %v", err)
	}
	if val2 != "secret-value" {
		t.Errorf("expected cached value %q, got %q", "secret-value", val2)
	}
	if stub.calls != 1 {
		t.Errorf("expected still 1 fetch call after cache hit, got %d", stub.calls)
	}
}

func TestLoad_PropagatesFetchError(t *testing.T) {
	fetchErr := errors.New("vault unavailable")
	stub := &stubFetch{}
	stub.add("", fetchErr)

	cp := checkpoint.New(5, 30*time.Second)
	_, err := cp.Load(context.Background(), "badkey", stub.fetch)
	if !errors.Is(err, fetchErr) {
		t.Errorf("expected fetch error, got %v", err)
	}
}

func TestInvalidate_ForcesFetch(t *testing.T) {
	stub := &stubFetch{}
	stub.add("first", nil)
	stub.add("second", nil)

	cp := checkpoint.New(5, 30*time.Second)
	_, _ = cp.Load(context.Background(), "k", stub.fetch)
	cp.Invalidate("k")

	val, err := cp.Load(context.Background(), "k", stub.fetch)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "second" {
		t.Errorf("expected refreshed value %q, got %q", "second", val)
	}
	if stub.calls != 2 {
		t.Errorf("expected 2 fetch calls after invalidation, got %d", stub.calls)
	}
}

func TestInvalidateAll_ClearsAllKeys(t *testing.T) {
	stub := &stubFetch{}
	stub.add("a1", nil)
	stub.add("b1", nil)
	stub.add("a2", nil)
	stub.add("b2", nil)

	cp := checkpoint.New(10, 30*time.Second)
	_, _ = cp.Load(context.Background(), "a", stub.fetch)
	_, _ = cp.Load(context.Background(), "b", stub.fetch)
	cp.InvalidateAll()

	cp.Load(context.Background(), "a", stub.fetch) //nolint:errcheck
	cp.Load(context.Background(), "b", stub.fetch) //nolint:errcheck

	if stub.calls != 4 {
		t.Errorf("expected 4 total fetch calls, got %d", stub.calls)
	}
}

func TestLoad_ContextCancelled_ReturnsError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	blocking := func(ctx context.Context, key string) (string, error) {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(5 * time.Second):
			return "late", nil
		}
	}

	cp := checkpoint.New(5, 30*time.Second)
	_, err := cp.Load(ctx, "key", blocking)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}
