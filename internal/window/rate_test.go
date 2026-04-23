package window

import (
	"strings"
	"testing"
	"time"
)

func TestPerSecond_ZeroWhenEmpty(t *testing.T) {
	r := NewRate(time.Second, 4)
	if r.PerSecond() != 0 {
		t.Fatal("expected 0 per second on empty window")
	}
}

func TestPerSecond_Calculation(t *testing.T) {
	r := NewRate(10*time.Second, 1)
	r.Observe(100)
	// 100 events over 10 s = 10 events/s
	got := r.PerSecond()
	if got != 10.0 {
		t.Fatalf("expected 10.0, got %f", got)
	}
}

func TestExceeds_BelowLimit(t *testing.T) {
	r := NewRate(10*time.Second, 1)
	r.Observe(5) // 0.5/s
	if r.Exceeds(1.0) {
		t.Fatal("should not exceed 1.0 rps")
	}
}

func TestExceeds_AboveLimit(t *testing.T) {
	r := NewRate(time.Second, 1)
	r.Observe(50)
	if !r.Exceeds(10.0) {
		t.Fatal("should exceed 10.0 rps")
	}
}

func TestString_ContainsWindow(t *testing.T) {
	r := NewRate(5*time.Second, 2)
	s := r.String()
	if !strings.Contains(s, "5s") {
		t.Fatalf("expected window in string, got %q", s)
	}
}

func TestRate_ZeroSizeWindow(t *testing.T) {
	r := &Rate{w: New(0, 1)}
	if r.PerSecond() != 0 {
		t.Fatal("expected 0 for zero-size window")
	}
}
