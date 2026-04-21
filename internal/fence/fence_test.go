package fence_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/fence"
)

func TestAcquire_SucceedsWhenFree(t *testing.T) {
	f := fence.New()
	release, err := f.Acquire(context.Background(), "key1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer release()
	if f.Active() != 1 {
		t.Fatalf("expected 1 active, got %d", f.Active())
	}
}

func TestAcquire_FailsWhenHeld(t *testing.T) {
	f := fence.New()
	release, err := f.Acquire(context.Background(), "key1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer release()

	_, err2 := f.Acquire(context.Background(), "key1")
	if err2 != fence.ErrFenced {
		t.Fatalf("expected ErrFenced, got %v", err2)
	}
}

func TestRelease_FreesKey(t *testing.T) {
	f := fence.New()
	release, _ := f.Acquire(context.Background(), "mykey")
	release()
	if f.Active() != 0 {
		t.Fatalf("expected 0 active after release, got %d", f.Active())
	}
}

func TestWait_BlocksUntilRelease(t *testing.T) {
	f := fence.New()
	release, _ := f.Acquire(context.Background(), "slow")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := f.Wait(context.Background(), "slow"); err != nil {
			t.Errorf("Wait returned unexpected error: %v", err)
		}
	}()

	time.Sleep(20 * time.Millisecond)
	release()
	wg.Wait()
}

func TestWait_ReturnsImmediatelyWhenFree(t *testing.T) {
	f := fence.New()
	ctx := context.Background()
	if err := f.Wait(ctx, "absent"); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestWait_RespectsContextCancellation(t *testing.T) {
	f := fence.New()
	release, _ := f.Acquire(context.Background(), "blocked")
	defer release()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	err := f.Wait(ctx, "blocked")
	if err == nil {
		t.Fatal("expected context error, got nil")
	}
}

func TestActive_MultipleKeys(t *testing.T) {
	f := fence.New()
	r1, _ := f.Acquire(context.Background(), "a")
	r2, _ := f.Acquire(context.Background(), "b")
	defer r1()
	defer r2()
	if f.Active() != 2 {
		t.Fatalf("expected 2, got %d", f.Active())
	}
}
