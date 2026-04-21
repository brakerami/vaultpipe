package semaphore_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/semaphore"
)

func TestNew_InvalidCap(t *testing.T) {
	_, err := semaphore.New(0)
	if err == nil {
		t.Fatal("expected error for zero capacity")
	}
}

func TestNew_ValidCap(t *testing.T) {
	sem, err := semaphore.New(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sem.Cap() != 3 {
		t.Fatalf("expected cap 3, got %d", sem.Cap())
	}
	if sem.Available() != 3 {
		t.Fatalf("expected 3 available, got %d", sem.Available())
	}
}

func TestAcquire_Release(t *testing.T) {
	sem, _ := semaphore.New(2)
	ctx := context.Background()

	if err := sem.Acquire(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sem.Available() != 1 {
		t.Fatalf("expected 1 available after acquire, got %d", sem.Available())
	}
	sem.Release()
	if sem.Available() != 2 {
		t.Fatalf("expected 2 available after release, got %d", sem.Available())
	}
}

func TestTryAcquire_FailsWhenFull(t *testing.T) {
	sem, _ := semaphore.New(1)
	if !sem.TryAcquire() {
		t.Fatal("expected TryAcquire to succeed on empty semaphore")
	}
	if sem.TryAcquire() {
		t.Fatal("expected TryAcquire to fail when semaphore is full")
	}
	sem.Release()
}

func TestAcquire_BlocksAtCap(t *testing.T) {
	sem, _ := semaphore.New(1)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_ = sem.TryAcquire() // fill the slot

	err := sem.Acquire(ctx)
	if err == nil {
		t.Fatal("expected context deadline error")
	}
}

func TestAcquire_UnblocksAfterRelease(t *testing.T) {
	sem, _ := semaphore.New(1)
	ctx := context.Background()

	_ = sem.TryAcquire()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := sem.Acquire(ctx); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		sem.Release()
	}()

	time.Sleep(10 * time.Millisecond)
	sem.Release()
	wg.Wait()
}

func TestRelease_PanicsWithoutAcquire(t *testing.T) {
	sem, _ := semaphore.New(1)
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on Release without Acquire")
		}
	}()
	sem.Release()
}
