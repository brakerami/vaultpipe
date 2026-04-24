package suppress

import (
	"context"
	"errors"
	"testing"
	"time"
)

var errVault = errors.New("vault unavailable")

func okFetch(_ context.Context, _ string) (string, error) {
	return "s3cr3t", nil
}

func errFetch(_ context.Context, _ string) (string, error) {
	return "", errVault
}

func TestGuardedFetch_SuccessPassthrough(t *testing.T) {
	s := New(time.Minute)
	fetch := GuardedFetch(s, okFetch)
	val, err := fetch(context.Background(), "secret/db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "s3cr3t" {
		t.Fatalf("expected s3cr3t, got %q", val)
	}
}

func TestGuardedFetch_FirstError_Forwarded(t *testing.T) {
	s := New(time.Minute)
	fetch := GuardedFetch(s, errFetch)
	_, err := fetch(context.Background(), "secret/db")
	if !errors.Is(err, errVault) {
		t.Fatalf("expected errVault, got %v", err)
	}
}

func TestGuardedFetch_SecondError_Suppressed(t *testing.T) {
	s := New(time.Minute)
	fetch := GuardedFetch(s, errFetch)
	fetch(context.Background(), "secret/db") // first: forwarded
	_, err := fetch(context.Background(), "secret/db") // second: suppressed
	if err != nil {
		t.Fatalf("expected suppressed nil error, got %v", err)
	}
}

func TestGuardedFetch_ErrorResetsAfterWindow(t *testing.T) {
	s := New(50 * time.Millisecond)
	now := time.Now()
	s.nowFn = func() time.Time { return now }

	fetch := GuardedFetch(s, errFetch)
	fetch(context.Background(), "secret/db")

	s.nowFn = func() time.Time { return now.Add(100 * time.Millisecond) }
	_, err := fetch(context.Background(), "secret/db")
	if !errors.Is(err, errVault) {
		t.Fatalf("expected error after window reset, got %v", err)
	}
}

func TestGuardedFetch_IndependentPaths(t *testing.T) {
	s := New(time.Minute)
	fetch := GuardedFetch(s, errFetch)
	_, err1 := fetch(context.Background(), "secret/db")
	_, err2 := fetch(context.Background(), "secret/api")
	if !errors.Is(err1, errVault) {
		t.Fatalf("expected errVault for db, got %v", err1)
	}
	if !errors.Is(err2, errVault) {
		t.Fatalf("expected errVault for api, got %v", err2)
	}
}
