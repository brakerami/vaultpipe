package cooldown

import (
	"context"
	"errors"
	"testing"
	"time"
)

func okFetch(_ context.Context, _ string) (string, error) {
	return "s3cr3t", nil
}

func errFetch(_ context.Context, _ string) (string, error) {
	return "", errors.New("vault unavailable")
}

func TestGuardedFetch_FirstCall_Succeeds(t *testing.T) {
	c := New(time.Minute)
	gf := GuardedFetch(c, okFetch)

	val, err := gf(context.Background(), "secret/db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "s3cr3t" {
		t.Fatalf("unexpected value: %q", val)
	}
}

func TestGuardedFetch_SecondCallBlocked(t *testing.T) {
	c := New(time.Minute)
	gf := GuardedFetch(c, okFetch)

	gf(context.Background(), "secret/db") //nolint:errcheck
	_, err := gf(context.Background(), "secret/db")
	if err == nil {
		t.Fatal("expected cooldown error on second call")
	}
	var ce *ErrCooldown
	if !errors.As(err, &ce) {
		t.Fatalf("expected *ErrCooldown, got %T", err)
	}
	if ce.Path != "secret/db" {
		t.Fatalf("unexpected path in error: %q", ce.Path)
	}
}

func TestGuardedFetch_ErrorResetsKey(t *testing.T) {
	c := New(time.Minute)
	gf := GuardedFetch(c, errFetch)

	// First call fails — key should be reset so a retry is allowed.
	gf(context.Background(), "secret/db") //nolint:errcheck

	// Replace with a successful fetch to confirm the key was reset.
	gf2 := GuardedFetch(c, okFetch)
	_, err := gf2(context.Background(), "secret/db")
	if err != nil {
		t.Fatalf("expected retry to succeed after error reset, got: %v", err)
	}
}

func TestGuardedFetch_IndependentPaths(t *testing.T) {
	c := New(time.Minute)
	gf := GuardedFetch(c, okFetch)

	gf(context.Background(), "secret/a") //nolint:errcheck
	_, err := gf(context.Background(), "secret/b")
	if err != nil {
		t.Fatalf("different path should not be throttled: %v", err)
	}
}

func TestErrCooldown_ErrorString(t *testing.T) {
	e := &ErrCooldown{Path: "secret/x", Remaining: 3 * time.Second}
	s := e.Error()
	if s == "" {
		t.Fatal("expected non-empty error string")
	}
}
