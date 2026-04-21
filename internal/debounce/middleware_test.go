package debounce_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/your-org/vaultpipe/internal/debounce"
)

func TestDebouncedRenewer_PanicsOnNilFunc(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil RenewFunc")
		}
	}()
	debounce.NewDebouncedRenewer(10*time.Millisecond, nil)
}

func TestDebouncedRenewer_CoalescesPerKey(t *testing.T) {
	var calls atomic.Int32
	dr := debounce.NewDebouncedRenewer(50*time.Millisecond, func(ctx context.Context, key string) error {
		calls.Add(1)
		return nil
	})

	for i := 0; i < 4; i++ {
		dr.Trigger(context.Background(), "db/creds")
		time.Sleep(10 * time.Millisecond)
	}
	time.Sleep(80 * time.Millisecond)

	if calls.Load() != 1 {
		t.Fatalf("expected 1 renewal, got %d", calls.Load())
	}
}

func TestDebouncedRenewer_IndependentKeys(t *testing.T) {
	var calls atomic.Int32
	dr := debounce.NewDebouncedRenewer(30*time.Millisecond, func(ctx context.Context, key string) error {
		calls.Add(1)
		return nil
	})

	dr.Trigger(context.Background(), "key-a")
	dr.Trigger(context.Background(), "key-b")
	time.Sleep(70 * time.Millisecond)

	if calls.Load() != 2 {
		t.Fatalf("expected 2 renewals (one per key), got %d", calls.Load())
	}
}

func TestDebouncedRenewer_Stop_CancelsAll(t *testing.T) {
	var calls atomic.Int32
	dr := debounce.NewDebouncedRenewer(60*time.Millisecond, func(ctx context.Context, key string) error {
		calls.Add(1)
		return nil
	})

	dr.Trigger(context.Background(), "key-x")
	dr.Trigger(context.Background(), "key-y")
	dr.Stop()
	time.Sleep(80 * time.Millisecond)

	if calls.Load() != 0 {
		t.Fatalf("expected 0 renewals after Stop, got %d", calls.Load())
	}
}
