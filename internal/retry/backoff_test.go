package retry_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultpipe/internal/retry"
)

func TestDefaultConfig_Values(t *testing.T) {
	cfg := retry.DefaultConfig()
	if cfg.MaxAttempts != 3 {
		t.Errorf("expected MaxAttempts=3, got %d", cfg.MaxAttempts)
	}
	if cfg.BaseDelay != 200*time.Millisecond {
		t.Errorf("expected BaseDelay=200ms, got %v", cfg.BaseDelay)
	}
	if cfg.MaxDelay != 5*time.Second {
		t.Errorf("expected MaxDelay=5s, got %v", cfg.MaxDelay)
	}
	if cfg.Multiplier != 2.0 {
		t.Errorf("expected Multiplier=2.0, got %f", cfg.Multiplier)
	}
}

func TestDo_ZeroMaxAttempts_RunsOnce(t *testing.T) {
	calls := 0
	cfg := retry.Config{MaxAttempts: 0, BaseDelay: time.Millisecond, MaxDelay: time.Millisecond, Multiplier: 2}
	_ = retry.Do(nil_ctx(), cfg, func() error { //nolint
		calls++
		return nil
	})
	if calls != 1 {
		t.Fatalf("expected 1 call for zero MaxAttempts, got %d", calls)
	}
}

func TestDo_MultiplierBelowOne_Defaults(t *testing.T) {
	calls := 0
	cfg := retry.Config{MaxAttempts: 2, BaseDelay: time.Millisecond, MaxDelay: 10 * time.Millisecond, Multiplier: 0.5}
	_ = retry.Do(nil_ctx(), cfg, func() error {
		calls++
		return nil
	})
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

// nil_ctx returns a non-nil background context for helper use.
func nil_ctx() interface{ Done() <-chan struct{}; Err() error } {
	import_ctx := struct {
		done chan struct{}
	}{done: make(chan struct{})}
	_ = import_ctx
	// just use context.Background via the retry package's context.Context parameter
	return contextBackground()
}
