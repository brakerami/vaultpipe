package jitter_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/jitter"
)

// fixedSource always returns the same value, making output deterministic.
type fixedSource struct{ v float64 }

func (f fixedSource) Float64() float64 { return f.v }

func TestJitterWith_ZeroFactor(t *testing.T) {
	base := 10 * time.Second
	got := jitter.JitterWith(base, 0, fixedSource{0.9})
	if got != base {
		t.Fatalf("expected %v, got %v", base, got)
	}
}

func TestJitterWith_FullFactor(t *testing.T) {
	base := 10 * time.Second
	// factor=1, rand=0.5 → base + 0.5*base = 15s
	got := jitter.JitterWith(base, 1.0, fixedSource{0.5})
	want := 15 * time.Second
	if got != want {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestJitterWith_ClampsFactor(t *testing.T) {
	base := 10 * time.Second
	// factor > 1 should be clamped to 1; rand=1.0 → base*2
	got := jitter.JitterWith(base, 5.0, fixedSource{1.0})
	want := 20 * time.Second
	if got != want {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestJitterWith_NegativeFactor(t *testing.T) {
	base := 10 * time.Second
	got := jitter.JitterWith(base, -1.0, fixedSource{0.5})
	if got != base {
		t.Fatalf("expected %v, got %v", base, got)
	}
}

func TestJitterWith_ZeroBase(t *testing.T) {
	got := jitter.JitterWith(0, 0.5, fixedSource{0.9})
	if got != 0 {
		t.Fatalf("expected 0, got %v", got)
	}
}

func TestFullWith_ReturnsWithinRange(t *testing.T) {
	base := 20 * time.Second
	// rand=0.75 → 15s
	got := jitter.FullWith(base, fixedSource{0.75})
	want := 15 * time.Second
	if got != want {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestFullWith_ZeroBase(t *testing.T) {
	got := jitter.FullWith(0, fixedSource{0.5})
	if got != 0 {
		t.Fatalf("expected 0, got %v", got)
	}
}

func TestEqualWith_ReturnsAtLeastHalf(t *testing.T) {
	base := 20 * time.Second
	// half=10s, rand=0 → 10s + 0 = 10s
	got := jitter.EqualWith(base, fixedSource{0.0})
	want := 10 * time.Second
	if got != want {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestEqualWith_ReturnsAtMostBase(t *testing.T) {
	base := 20 * time.Second
	// half=10s, rand=1 → 10s + 10s = 20s
	got := jitter.EqualWith(base, fixedSource{1.0})
	if got > base {
		t.Fatalf("expected at most %v, got %v", base, got)
	}
}
