package limiter_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"vaultpipe/internal/limiter"
)

func TestNew_InvalidCap(t *testing.T) {
	_, err := limiter.New(0)
	if err == nil {
		t.Fatal("expected error for cap=0, got nil")
	}
}

func TestNew_ValidCap(t *testing.T) {
	l, err := limiter.New(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.Cap() != 3 {
		t.Fatalf("expected cap 3, got %d", l.Cap())
	}
}

func TestAcquire_Release(t *testing.T) {
	l, _ := limiter.New(2)
	ctx := context.Background()

	if err := l.Acquire(ctx); err != nil {
		t.Fatalf("first acquire failed: %v", err)
	}
	if l.InFlight() != 1 {
		t.Fatalf("expected 1 in-flight, got %d", l.InFlight())
	}
	l.Release()
	if l.InFlight() != 0 {
		t.Fatalf("expected 0 in-flight after release, got %d", l.InFlight())
	}
}

func TestAcquire_BlocksAtCap(t *testing.T) {
	l, _ := limiter.New(1)
	ctx := context.Background()

	if err := l.Acquire(ctx); err != nil {
		t.Fatalf("first acquire: %v", err)
	}

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := l.Acquire(ctxTimeout)
	if err == nil {
		t.Fatal("expected error when limit is full, got nil")
	}
}

func TestAcquire_UnblocksAfterRelease(t *testing.T) {
	l, _ := limiter.New(1)
	ctx := context.Background()

	if err := l.Acquire(ctx); err != nil {
		t.Fatalf("first acquire: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := l.Acquire(ctx); err != nil {
			t.Errorf("second acquire failed: %v", err)
		}
		l.Release()
	}()

	time.Sleep(20 * time.Millisecond)
	l.Release()
	wg.Wait()
}

func TestInFlight_Concurrent(t *testing.T) {
	l, _ := limiter.New(5)
	ctx := context.Background()
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = l.Acquire(ctx)
		}()
	}
	wg.Wait()

	if l.InFlight() != 5 {
		t.Fatalf("expected 5 in-flight, got %d", l.InFlight())
	}
}
