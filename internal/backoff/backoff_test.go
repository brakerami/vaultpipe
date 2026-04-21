package backoff_test

import (
	"testing"
	"time"

	"vaultpipe/internal/backoff"
)

func TestDefault_ReturnsReasonableValues(t *testing.T) {
	cfg := backoff.Default()
	if cfg.Base == 0 {
		t.Fatal("expected non-zero Base")
	}
	if cfg.Max == 0 {
		t.Fatal("expected non-zero Max")
	}
	if cfg.Multiplier < 1.0 {
		t.Fatalf("expected Multiplier >= 1, got %v", cfg.Multiplier)
	}
}

func TestNext_FirstAttempt_ReturnsBase(t *testing.T) {
	cfg := backoff.Config{
		Base:       100 * time.Millisecond,
		Max:        10 * time.Second,
		Multiplier: 2.0,
		Jitter:     0,
	}
	got := cfg.Next(0)
	if got != 100*time.Millisecond {
		t.Fatalf("expected 100ms, got %v", got)
	}
}

func TestNext_GrowsExponentially(t *testing.T) {
	cfg := backoff.Config{
		Base:       100 * time.Millisecond,
		Max:        10 * time.Second,
		Multiplier: 2.0,
		Jitter:     0,
	}
	prev := cfg.Next(0)
	for i := 1; i <= 4; i++ {
		next := cfg.Next(i)
		if next <= prev {
			t.Fatalf("attempt %d: expected delay > %v, got %v", i, prev, next)
		}
		prev = next
	}
}

func TestNext_CapsAtMax(t *testing.T) {
	cfg := backoff.Config{
		Base:       1 * time.Second,
		Max:        2 * time.Second,
		Multiplier: 10.0,
		Jitter:     0,
	}
	for i := 0; i < 10; i++ {
		d := cfg.Next(i)
		// allow up to Max (jitter is 0 so exact cap applies)
		if d > cfg.Max {
			t.Fatalf("attempt %d: delay %v exceeds max %v", i, d, cfg.Max)
		}
	}
}

func TestNext_JitterAddsVariance(t *testing.T) {
	cfg := backoff.Config{
		Base:       100 * time.Millisecond,
		Max:        10 * time.Second,
		Multiplier: 2.0,
		Jitter:     0.5,
	}
	// Run enough iterations to observe that not all values are identical.
	seen := map[time.Duration]bool{}
	for i := 0; i < 20; i++ {
		seen[cfg.Next(1)] = true
	}
	if len(seen) < 2 {
		t.Fatal("expected jitter to produce varied delays")
	}
}

func TestNext_BadMultiplier_Defaults(t *testing.T) {
	cfg := backoff.Config{
		Base:       100 * time.Millisecond,
		Max:        10 * time.Second,
		Multiplier: 0.5, // invalid — should default to 2.0
		Jitter:     0,
	}
	d0 := cfg.Next(0)
	d1 := cfg.Next(1)
	if d1 <= d0 {
		t.Fatalf("expected growth with corrected multiplier, got d0=%v d1=%v", d0, d1)
	}
}
