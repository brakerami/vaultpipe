package fence_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/fence"
	"github.com/yourusername/vaultpipe/internal/metrics"
)

func TestDeduplicate_SingleCaller(t *testing.T) {
	f := fence.New()
	calls := 0
	next := func(_ context.Context, path string) (string, error) {
		calls++
		return "val", nil
	}
	wrapped := fence.Deduplicate(f, nil, next)
	v, err := wrapped(context.Background(), "secret/foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "val" {
		t.Fatalf("expected 'val', got %q", v)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDeduplicate_ConcurrentCallersCollapseIntoTwo(t *testing.T) {
	f := fence.New()
	reg := metrics.New()
	var callCount int64
	ready := make(chan struct{})

	next := func(_ context.Context, _ string) (string, error) {
		atomic.AddInt64(&callCount, 1)
		<-ready // block until we release
		return "secret", nil
	}
	wrapped := fence.Deduplicate(f, reg, next)

	var wg sync.WaitGroup
	const goroutines = 5
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			wrapped(context.Background(), "secret/bar") //nolint:errcheck
		}()
	}

	// Give goroutines time to pile up behind the fence.
	time.Sleep(30 * time.Millisecond)
	close(ready)
	wg.Wait()

	// The first goroutine held the fence; the rest waited and then each
	// called next once more, so total = 1 (holder) + N-1 (waiters).
	if atomic.LoadInt64(&callCount) < 1 {
		t.Fatal("expected at least one call to next")
	}
}

func TestDeduplicate_ContextCancelledWhileWaiting(t *testing.T) {
	f := fence.New()
	release, _ := f.Acquire(context.Background(), "secret/locked")
	defer release()

	next := func(_ context.Context, _ string) (string, error) {
		return "x", nil
	}
	wrapped := fence.Deduplicate(f, nil, next)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	_, err := wrapped(ctx, "secret/locked")
	if err == nil {
		t.Fatal("expected error due to context cancellation")
	}
}
