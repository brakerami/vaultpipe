package watch_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/your-org/vaultpipe/internal/watch"
)

func TestNew_DefaultThreshold(t *testing.T) {
	w := watch.New(func(_ context.Context, _ string) error { return nil }, 0)
	if w == nil {
		t.Fatal("expected non-nil watcher")
	}
}

func TestAdd_TriggersRenewal(t *testing.T) {
	var called atomic.Int32

	w := watch.New(func(_ context.Context, ref string) error {
		if ref == "secret/data/test" {
			called.Add(1)
		}
		return nil
	}, 0.9) // renew when 90% of TTL elapsed → fires almost immediately

	ctx := context.Background()
	w.Add(ctx, "secret/data/test", 50*time.Millisecond)

	time.Sleep(100 * time.Millisecond)

	if called.Load() == 0 {
		t.Error("expected renewal callback to be called")
	}
}

func TestRemove_StopsRenewal(t *testing.T) {
	var called atomic.Int32

	w := watch.New(func(_ context.Context, _ string) error {
		called.Add(1)
		return nil
	}, 0.1)

	ctx := context.Background()
	w.Add(ctx, "secret/data/gone", 200*time.Millisecond)
	w.Remove("secret/data/gone")

	time.Sleep(250 * time.Millisecond)

	if called.Load() != 0 {
		t.Error("expected no renewal after Remove")
	}
}

func TestStop_CancelsAll(t *testing.T) {
	var called atomic.Int32

	w := watch.New(func(_ context.Context, _ string) error {
		called.Add(1)
		return nil
	}, 0.1)

	ctx := context.Background()
	w.Add(ctx, "secret/a", 300*time.Millisecond)
	w.Add(ctx, "secret/b", 300*time.Millisecond)
	w.Stop()

	time.Sleep(350 * time.Millisecond)

	if called.Load() != 0 {
		t.Errorf("expected 0 renewals after Stop, got %d", called.Load())
	}
}

func TestAdd_ResetExistingLease(t *testing.T) {
	var called atomic.Int32

	w := watch.New(func(_ context.Context, _ string) error {
		called.Add(1)
		return nil
	}, 0.9)

	ctx := context.Background()
	// Add once with very short TTL, then immediately re-add with longer TTL.
	w.Add(ctx, "secret/reset", 20*time.Millisecond)
	w.Add(ctx, "secret/reset", 500*time.Millisecond)

	time.Sleep(80 * time.Millisecond)

	// First watcher was cancelled; second hasn't fired yet.
	if called.Load() != 0 {
		t.Errorf("expected 0 calls after lease reset, got %d", called.Load())
	}
	w.Stop()
}
