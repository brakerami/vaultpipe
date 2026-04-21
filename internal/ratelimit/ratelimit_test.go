package ratelimit_test

import (
	"context"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/ratelimit"
)

func TestNew_InvalidRPS(t *testing.T) {
	_, err := ratelimit.New(0)
	if err == nil {
		t.Fatal("expected error for rps=0, got nil")
	}
	_, err = ratelimit.New(-5)
	if err == nil {
		t.Fatal("expected error for negative rps, got nil")
	}
}

func TestNew_ValidRPS(t *testing.T) {
	l, err := ratelimit.New(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l == nil {
		t.Fatal("expected non-nil limiter")
	}
}

func TestWait_AcquiresToken(t *testing.T) {
	l, _ := ratelimit.New(100)
	ctx := context.Background()
	if err := l.Wait(ctx); err != nil {
		t.Fatalf("unexpected error on first Wait: %v", err)
	}
}

func TestWait_CancelledContext(t *testing.T) {
	// Very low rate so no token is immediately available after the first burst.
	l, _ := ratelimit.New(0.001)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	// Drain the initial token.
	_ = l.Wait(context.Background())
	// Now the bucket is empty; next Wait should respect context cancellation.
	err := l.Wait(ctx)
	if err == nil {
		t.Fatal("expected context deadline error, got nil")
	}
}

func TestWait_MultipleTokens(t *testing.T) {
	l, _ := ratelimit.New(50)
	ctx := context.Background()
	for i := 0; i < 10; i++ {
		if err := l.Wait(ctx); err != nil {
			t.Fatalf("Wait %d failed: %v", i, err)
		}
	}
}

func TestAvailable_DecreasesAfterWait(t *testing.T) {
	l, _ := ratelimit.New(10)
	before := l.Available()
	_ = l.Wait(context.Background())
	after := l.Available()
	if after >= before {
		t.Errorf("expected Available to decrease after Wait: before=%v after=%v", before, after)
	}
}
