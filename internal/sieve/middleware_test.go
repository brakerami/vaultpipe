package sieve_test

import (
	"context"
	"errors"
	"testing"

	"github.com/yourusername/vaultpipe/internal/sieve"
)

func okFetch(_ context.Context, _ string) (map[string]string, error) {
	return map[string]string{"key": "value"}, nil
}

func errFetch(_ context.Context, _ string) (map[string]string, error) {
	return nil, errors.New("vault: not found")
}

func TestGuardedFetch_AllowedPath_CallsNext(t *testing.T) {
	s := sieve.New()
	s.Allow("secret/")
	guarded := sieve.GuardedFetch(s, okFetch)
	vals, err := guarded(context.Background(), "secret/data/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vals["key"] != "value" {
		t.Fatalf("expected value 'value', got %q", vals["key"])
	}
}

func TestGuardedFetch_DeniedPath_ReturnsError(t *testing.T) {
	s := sieve.New()
	s.Deny("secret/data/internal")
	guarded := sieve.GuardedFetch(s, okFetch)
	_, err := guarded(context.Background(), "secret/data/internal/creds")
	if err == nil {
		t.Fatal("expected error for denied path")
	}
}

func TestGuardedFetch_PropagatesNextError(t *testing.T) {
	s := sieve.New()
	guarded := sieve.GuardedFetch(s, errFetch)
	_, err := guarded(context.Background(), "secret/data/app")
	if err == nil {
		t.Fatal("expected error from next fetch")
	}
}

func TestGuardedFetch_NilSieve_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil Sieve")
		}
	}()
	sieve.GuardedFetch(nil, okFetch)
}

func TestGuardedFetch_NilNext_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil FetchFunc")
		}
	}()
	sieve.GuardedFetch(sieve.New(), nil)
}
