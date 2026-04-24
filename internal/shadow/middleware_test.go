package shadow_test

import (
	"context"
	"errors"
	"testing"

	"github.com/yourusername/vaultpipe/internal/shadow"
)

func makeFetch(val string, err error) shadow.FetchFunc {
	return func(_ context.Context, _ string) (string, error) {
		return val, err
	}
}

func TestTrackingFetch_RecordsValue(t *testing.T) {
	s := shadow.New()
	fetch := shadow.TrackingFetch(s, nil, makeFetch("v1", nil))

	v, err := fetch(context.Background(), "secret/db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "v1" {
		t.Errorf("got %q, want %q", v, "v1")
	}
	e, ok := s.Get("secret/db")
	if !ok || e.Value != "v1" {
		t.Error("expected shadow to store fetched value")
	}
}

func TestTrackingFetch_CallsOnDrift(t *testing.T) {
	s := shadow.New()
	s.Set("secret/db", "old")

	var driftKey, driftOld, driftNew string
	onDrift := func(k, o, n string) { driftKey, driftOld, driftNew = k, o, n }

	fetch := shadow.TrackingFetch(s, onDrift, makeFetch("new", nil))
	_, _ = fetch(context.Background(), "secret/db")

	if driftKey != "secret/db" {
		t.Errorf("expected drift key %q, got %q", "secret/db", driftKey)
	}
	if driftOld != "old" || driftNew != "new" {
		t.Errorf("unexpected drift values old=%q new=%q", driftOld, driftNew)
	}
}

func TestTrackingFetch_PropagatesError(t *testing.T) {
	s := shadow.New()
	fetch := shadow.TrackingFetch(s, nil, makeFetch("", errors.New("vault down")))

	_, err := fetch(context.Background(), "secret/db")
	if err == nil {
		t.Fatal("expected error to propagate")
	}
}

func TestStrictFetch_NoDrift_ReturnsValue(t *testing.T) {
	s := shadow.New()
	s.Set("secret/api", "stable")

	fetch := shadow.StrictFetch(s, makeFetch("stable", nil))
	v, err := fetch(context.Background(), "secret/api")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "stable" {
		t.Errorf("got %q, want %q", v, "stable")
	}
}

func TestStrictFetch_DriftDetected_ReturnsError(t *testing.T) {
	s := shadow.New()
	s.Set("secret/api", "old")

	fetch := shadow.StrictFetch(s, makeFetch("new", nil))
	_, err := fetch(context.Background(), "secret/api")

	var driftErr *shadow.ErrDriftDetected
	if !errors.As(err, &driftErr) {
		t.Fatalf("expected ErrDriftDetected, got %v", err)
	}
	if driftErr.Key != "secret/api" {
		t.Errorf("unexpected drift key: %q", driftErr.Key)
	}
	// shadow must NOT be updated on drift
	e, _ := s.Get("secret/api")
	if e.Value != "old" {
		t.Error("shadow should remain unchanged after drift error")
	}
}
