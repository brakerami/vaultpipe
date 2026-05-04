package latch_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"vaultpipe/internal/latch"
)

func TestNew_NotSetInitially(t *testing.T) {
	l := latch.New()
	if l.IsSet() {
		t.Fatal("expected latch to be unset initially")
	}
}

func TestSet_ClosesLatch(t *testing.T) {
	l := latch.New()
	l.Set()
	if !l.IsSet() {
		t.Fatal("expected latch to be set after Set()")
	}
}

func TestSet_IdempotentMultipleCalls(t *testing.T) {
	l := latch.New()
	for i := 0; i < 10; i++ {
		l.Set() // must not panic
	}
	if !l.IsSet() {
		t.Fatal("latch should remain set")
	}
}

func TestWait_ReturnsWhenSet(t *testing.T) {
	l := latch.New()
	go func() {
		time.Sleep(20 * time.Millisecond)
		l.Set()
	}()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := l.Wait(ctx); err != nil {
		t.Fatalf("Wait returned error: %v", err)
	}
}

func TestWait_CancelledContext(t *testing.T) {
	l := latch.New()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := l.Wait(ctx); err == nil {
		t.Fatal("expected context cancellation error")
	}
}

func TestC_ClosedAfterSet(t *testing.T) {
	l := latch.New()
	l.Set()
	select {
	case <-l.C():
		// ok
	default:
		t.Fatal("channel should be closed after Set")
	}
}

func TestSet_ConcurrentSafety(t *testing.T) {
	l := latch.New()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			l.Set()
		}()
	}
	wg.Wait()
	if !l.IsSet() {
		t.Fatal("latch must be set after concurrent calls")
	}
}
