package debounce_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/your-org/vaultpipe/internal/debounce"
)

func TestNew_PanicsOnZeroWait(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for zero wait")
		}
	}()
	debounce.New(0, func(ctx context.Context) {})
}

func TestNew_PanicsOnNilFunc(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil fn")
		}
	}()
	debounce.New(10*time.Millisecond, nil)
}

func TestTrigger_InvokesFnAfterWait(t *testing.T) {
	var called atomic.Int32
	d := debounce.New(30*time.Millisecond, func(ctx context.Context) {
		called.Add(1)
	})
	d.Trigger(context.Background())
	time.Sleep(60 * time.Millisecond)
	if called.Load() != 1 {
		t.Fatalf("expected fn called once, got %d", called.Load())
	}
}

func TestTrigger_CoalescesRapidCalls(t *testing.T) {
	var called atomic.Int32
	d := debounce.New(50*time.Millisecond, func(ctx context.Context) {
		called.Add(1)
	})
	for i := 0; i < 5; i++ {
		d.Trigger(context.Background())
		time.Sleep(10 * time.Millisecond)
	}
	time.Sleep(80 * time.Millisecond)
	if called.Load() != 1 {
		t.Fatalf("expected exactly 1 call, got %d", called.Load())
	}
}

func TestStop_CancelsPendingInvocation(t *testing.T) {
	var called atomic.Int32
	d := debounce.New(40*time.Millisecond, func(ctx context.Context) {
		called.Add(1)
	})
	d.Trigger(context.Background())
	d.Stop()
	time.Sleep(60 * time.Millisecond)
	if called.Load() != 0 {
		t.Fatalf("expected fn not called after Stop, got %d", called.Load())
	}
}

func TestFlush_InvokesImmediately(t *testing.T) {
	var called atomic.Int32
	d := debounce.New(200*time.Millisecond, func(ctx context.Context) {
		called.Add(1)
	})
	d.Trigger(context.Background())
	d.Flush(context.Background())
	time.Sleep(20 * time.Millisecond)
	if called.Load() != 1 {
		t.Fatalf("expected fn called once after Flush, got %d", called.Load())
	}
}

func TestFlush_NoopWhenNoPending(t *testing.T) {
	var called atomic.Int32
	d := debounce.New(30*time.Millisecond, func(ctx context.Context) {
		called.Add(1)
	})
	d.Flush(context.Background()) // no pending — should not panic or call fn
	time.Sleep(10 * time.Millisecond)
	if called.Load() != 0 {
		t.Fatalf("expected fn not called, got %d", called.Load())
	}
}
