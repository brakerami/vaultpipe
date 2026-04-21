package rotate_test

import (
	"context"
	"sync"
	"testing"

	"github.com/yourusername/vaultpipe/internal/rotate"
)

func TestObserve_FirstObservation_NoHook(t *testing.T) {
	d := rotate.New()
	called := false
	d.OnChange(func(_ context.Context, _, _, _ string) { called = true })

	rotated := d.Observe(context.Background(), "DB_PASS", "secret1")
	if rotated {
		t.Error("expected false on first observation")
	}
	if called {
		t.Error("hook must not fire on first observation")
	}
}

func TestObserve_SameValue_NoHook(t *testing.T) {
	d := rotate.New()
	d.Observe(context.Background(), "DB_PASS", "secret1") //nolint:errcheck

	called := false
	d.OnChange(func(_ context.Context, _, _, _ string) { called = true })
	d.Observe(context.Background(), "DB_PASS", "secret1") //nolint:errcheck

	if called {
		t.Error("hook must not fire when value is unchanged")
	}
}

func TestObserve_ChangedValue_FiresHook(t *testing.T) {
	d := rotate.New()
	d.Observe(context.Background(), "API_KEY", "old") //nolint:errcheck

	var mu sync.Mutex
	var gotKey, gotOld, gotNew string
	d.OnChange(func(_ context.Context, key, oldV, newV string) {
		mu.Lock()
		gotKey, gotOld, gotNew = key, oldV, newV
		mu.Unlock()
	})

	rotated := d.Observe(context.Background(), "API_KEY", "new")
	if !rotated {
		t.Error("expected rotation to be detected")
	}
	mu.Lock()
	defer mu.Unlock()
	if gotKey != "API_KEY" || gotOld != "old" || gotNew != "new" {
		t.Errorf("unexpected hook args: key=%q old=%q new=%q", gotKey, gotOld, gotNew)
	}
}

func TestSeed_PreventsTriggerOnFirstObserve(t *testing.T) {
	d := rotate.New()
	d.Seed("TOKEN", "seeded")

	called := false
	d.OnChange(func(_ context.Context, _, _, _ string) { called = true })

	// same value as seed — no rotation
	d.Observe(context.Background(), "TOKEN", "seeded") //nolint:errcheck
	if called {
		t.Error("hook must not fire when value matches seed")
	}
}

func TestReset_AllowsReSeed(t *testing.T) {
	d := rotate.New()
	d.Observe(context.Background(), "X", "v1") //nolint:errcheck
	d.Reset("X")

	called := false
	d.OnChange(func(_ context.Context, _, _, _ string) { called = true })

	// After reset, next observe is treated as first — no hook.
	d.Observe(context.Background(), "X", "v2") //nolint:errcheck
	if called {
		t.Error("hook must not fire after reset on first re-observation")
	}
}
