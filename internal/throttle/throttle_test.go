package throttle_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/your-org/vaultpipe/internal/throttle"
)

func TestNew_PanicsOnZeroInterval(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for zero interval")
		}
	}()
	throttle.New(0)
}

func TestNew_PanicsOnNegativeInterval(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for negative interval")
		}
	}()
	throttle.New(-time.Second)
}

func TestDo_ExecutesImmediatelyOnFirstCall(t *testing.T) {
	th := throttle.New(100 * time.Millisecond)
	var called int32
	err := th.Do(context.Background(), func() error {
		atomic.AddInt32(&called, 1)
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if atomic.LoadInt32(&called) != 1 {
		t.Fatal("fn was not called")
	}
}

func TestDo_WaitsForCooldown(t *testing.T) {
	th := throttle.New(80 * time.Millisecond)
	var calls int32
	fn := func() error { atomic.AddInt32(&calls, 1); return nil }

	_ = th.Do(context.Background(), fn)

	start := time.Now()
	_ = th.Do(context.Background(), fn)
	elapsed := time.Since(start)

	if elapsed < 60*time.Millisecond {
		t.Fatalf("second call returned too quickly: %v", elapsed)
	}
	if atomic.LoadInt32(&calls) != 2 {
		t.Fatalf("expected 2 calls, got %d", atomic.LoadInt32(&calls))
	}
}

func TestDo_ContextCancelledDuringWait(t *testing.T) {
	th := throttle.New(500 * time.Millisecond)
	_ = th.Do(context.Background(), func() error { return nil })

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	err := th.Do(ctx, func() error { return nil })
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected DeadlineExceeded, got %v", err)
	}
}

func TestDo_PropagatesFnError(t *testing.T) {
	th := throttle.New(10 * time.Millisecond)
	sentinel := errors.New("vault unavailable")
	err := th.Do(context.Background(), func() error { return sentinel })
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}

func TestReset_AllowsImmediateReexecution(t *testing.T) {
	th := throttle.New(500 * time.Millisecond)
	_ = th.Do(context.Background(), func() error { return nil })
	th.Reset()

	start := time.Now()
	_ = th.Do(context.Background(), func() error { return nil })
	if time.Since(start) > 20*time.Millisecond {
		t.Fatal("Reset did not allow immediate re-execution")
	}
}

func TestLastRun_ZeroBeforeFirstCall(t *testing.T) {
	th := throttle.New(time.Second)
	if !th.LastRun().IsZero() {
		t.Fatal("LastRun should be zero before first call")
	}
}

func TestLastRun_UpdatedAfterCall(t *testing.T) {
	th := throttle.New(time.Second)
	before := time.Now()
	_ = th.Do(context.Background(), func() error { return nil })
	if th.LastRun().Before(before) {
		t.Fatal("LastRun not updated after Do")
	}
}
