package stagger_test

import (
	"context"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/stagger"
)

func TestNew_PanicsOnZeroWindow(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for zero window")
		}
	}()
	stagger.New(0)
}

func TestNew_PanicsOnNegativeWindow(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for negative window")
		}
	}()
	stagger.New(-time.Second)
}

func TestDelay_WithinWindow(t *testing.T) {
	window := 100 * time.Millisecond
	s := stagger.New(window)
	for i := 0; i < 50; i++ {
		d := s.Delay()
		if d < 0 || d >= window {
			t.Fatalf("delay %v outside [0, %v)", d, window)
		}
	}
}

func TestWait_CompletesWithinWindow(t *testing.T) {
	s := stagger.New(50 * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	if err := s.Wait(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWait_ContextCancelled(t *testing.T) {
	s := stagger.New(10 * time.Second) // very long window
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately
	err := s.Wait(ctx)
	if err == nil {
		t.Fatal("expected context error, got nil")
	}
}

func TestDo_CallsFnAfterDelay(t *testing.T) {
	s := stagger.New(20 * time.Millisecond)
	called := false
	err := s.Do(context.Background(), func() error {
		called = true
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("fn was not called")
	}
}

func TestDo_SkipsFnOnCancelledContext(t *testing.T) {
	s := stagger.New(10 * time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	called := false
	err := s.Do(ctx, func() error {
		called = true
		return nil
	})
	if err == nil {
		t.Fatal("expected error from cancelled context")
	}
	if called {
		t.Fatal("fn should not have been called")
	}
}
