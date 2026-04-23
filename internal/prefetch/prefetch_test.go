package prefetch_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/your-org/vaultpipe/internal/prefetch"
)

func TestRegister_TriggersOnRenew(t *testing.T) {
	var mu sync.Mutex
	renewed := []string{}

	p := prefetch.New(
		func(_ context.Context, path string) (string, error) {
			return "secret-" + path, nil
		},
		func(path, _ string) {
			mu.Lock()
			renewed = append(renewed, path)
			mu.Unlock()
		},
	)

	p.Register("kv/data/db", 100*time.Millisecond, 0.5)
	time.Sleep(120 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(renewed) == 0 {
		t.Fatal("expected onRenew to be called")
	}
	if renewed[0] != "kv/data/db" {
		t.Fatalf("unexpected path: %s", renewed[0])
	}
}

func TestDeregister_CancelsPrefetch(t *testing.T) {
	called := false

	p := prefetch.New(
		func(_ context.Context, path string) (string, error) {
			return "v", nil
		},
		func(_, _ string) { called = true },
	)

	p.Register("kv/data/token", 200*time.Millisecond, 0.5)
	p.Deregister("kv/data/token")
	time.Sleep(220 * time.Millisecond)

	if called {
		t.Fatal("onRenew should not have been called after deregister")
	}
}

func TestRegister_ReplacesExisting(t *testing.T) {
	var mu sync.Mutex
	count := 0

	p := prefetch.New(
		func(_ context.Context, _ string) (string, error) { return "v", nil },
		func(_, _ string) {
			mu.Lock()
			count++
			mu.Unlock()
		},
	)

	// Register twice — only the second should fire.
	p.Register("kv/data/x", 500*time.Millisecond, 0.5)
	p.Register("kv/data/x", 80*time.Millisecond, 0.5)
	time.Sleep(120 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if count != 1 {
		t.Fatalf("expected 1 renewal, got %d", count)
	}
}

func TestStop_CancelsAll(t *testing.T) {
	called := false

	p := prefetch.New(
		func(_ context.Context, _ string) (string, error) { return "v", nil },
		func(_, _ string) { called = true },
	)

	p.Register("kv/data/a", 300*time.Millisecond, 0.5)
	p.Register("kv/data/b", 300*time.Millisecond, 0.5)
	p.Stop()
	time.Sleep(320 * time.Millisecond)

	if called {
		t.Fatal("onRenew should not be called after Stop")
	}
}

func TestNew_NilOnRenew_DoesNotPanic(t *testing.T) {
	p := prefetch.New(
		func(_ context.Context, _ string) (string, error) { return "v", nil },
		nil,
	)
	p.Register("kv/data/safe", 60*time.Millisecond, 0.5)
	time.Sleep(80 * time.Millisecond) // should not panic
}
