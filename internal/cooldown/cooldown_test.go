package cooldown

import (
	"testing"
	"time"
)

func TestNew_PanicsOnZeroInterval(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for zero interval")
		}
	}()
	New(0)
}

func TestNew_PanicsOnNegativeInterval(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for negative interval")
		}
	}()
	New(-time.Second)
}

func TestAllow_FirstCall_Succeeds(t *testing.T) {
	c := New(time.Minute)
	if !c.Allow("k") {
		t.Fatal("expected first Allow to return true")
	}
}

func TestAllow_SecondCallWithinCooldown_Denied(t *testing.T) {
	c := New(time.Minute)
	c.Allow("k")
	if c.Allow("k") {
		t.Fatal("expected second Allow within cooldown to return false")
	}
}

func TestAllow_AfterCooldownExpires_Succeeds(t *testing.T) {
	now := time.Now()
	c := New(time.Second)
	c.now = func() time.Time { return now }
	c.Allow("k")

	// Advance clock past the interval.
	c.now = func() time.Time { return now.Add(2 * time.Second) }
	if !c.Allow("k") {
		t.Fatal("expected Allow to succeed after cooldown expired")
	}
}

func TestAllow_IndependentKeys(t *testing.T) {
	c := New(time.Minute)
	c.Allow("a")
	if !c.Allow("b") {
		t.Fatal("expected Allow for different key to succeed")
	}
}

func TestReset_AllowsImmediateRetry(t *testing.T) {
	c := New(time.Minute)
	c.Allow("k")
	c.Reset("k")
	if !c.Allow("k") {
		t.Fatal("expected Allow to succeed after Reset")
	}
}

func TestRemaining_ZeroWhenUnknownKey(t *testing.T) {
	c := New(time.Minute)
	if r := c.Remaining("missing"); r != 0 {
		t.Fatalf("expected 0, got %v", r)
	}
}

func TestRemaining_PositiveWithinCooldown(t *testing.T) {
	now := time.Now()
	c := New(10 * time.Second)
	c.now = func() time.Time { return now }
	c.Allow("k")

	c.now = func() time.Time { return now.Add(3 * time.Second) }
	r := c.Remaining("k")
	if r <= 0 || r > 10*time.Second {
		t.Fatalf("unexpected remaining duration: %v", r)
	}
}

func TestRemaining_ZeroAfterExpiry(t *testing.T) {
	now := time.Now()
	c := New(5 * time.Second)
	c.now = func() time.Time { return now }
	c.Allow("k")

	c.now = func() time.Time { return now.Add(10 * time.Second) }
	if r := c.Remaining("k"); r != 0 {
		t.Fatalf("expected 0, got %v", r)
	}
}
