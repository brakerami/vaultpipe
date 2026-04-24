package suppress

import (
	"testing"
	"time"
)

func TestNew_PanicsOnZeroWindow(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on zero window")
		}
	}()
	New(0)
}

func TestNew_PanicsOnNegativeWindow(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on negative window")
		}
	}()
	New(-time.Second)
}

func TestAllow_FirstCall_Succeeds(t *testing.T) {
	s := New(time.Minute)
	if !s.Allow("vault/secret/db") {
		t.Fatal("expected first Allow to return true")
	}
}

func TestAllow_SecondCallWithinWindow_Suppressed(t *testing.T) {
	s := New(time.Minute)
	s.Allow("key")
	if s.Allow("key") {
		t.Fatal("expected second Allow within window to return false")
	}
}

func TestAllow_AfterWindowExpires_Allowed(t *testing.T) {
	s := New(50 * time.Millisecond)
	now := time.Now()
	s.nowFn = func() time.Time { return now }
	s.Allow("key")

	s.nowFn = func() time.Time { return now.Add(100 * time.Millisecond) }
	if !s.Allow("key") {
		t.Fatal("expected Allow after window expiry to return true")
	}
}

func TestAllow_IndependentKeys(t *testing.T) {
	s := New(time.Minute)
	if !s.Allow("key-a") {
		t.Fatal("expected key-a to be allowed")
	}
	if !s.Allow("key-b") {
		t.Fatal("expected key-b to be allowed independently")
	}
}

func TestReset_ClearsState(t *testing.T) {
	s := New(time.Minute)
	s.Allow("key")
	s.Reset()
	if !s.Allow("key") {
		t.Fatal("expected Allow after Reset to return true")
	}
}

func TestLen_TracksKeys(t *testing.T) {
	s := New(time.Minute)
	s.Allow("a")
	s.Allow("b")
	s.Allow("a") // duplicate, not added again but already counted
	if s.Len() != 2 {
		t.Fatalf("expected Len 2, got %d", s.Len())
	}
}
