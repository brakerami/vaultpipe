package snapshot_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/snapshot"
)

func TestTake_StoresValues(t *testing.T) {
	s := snapshot.Take(map[string]string{"FOO": "bar", "BAZ": "qux"})
	v, ok := s.Get("FOO")
	if !ok || v != "bar" {
		t.Fatalf("expected FOO=bar, got %q ok=%v", v, ok)
	}
}

func TestTake_SetsTimestamp(t *testing.T) {
	before := time.Now().UTC()
	s := snapshot.Take(map[string]string{"X": "1"})
	after := time.Now().UTC()
	if s.TakenAt().Before(before) || s.TakenAt().After(after) {
		t.Fatalf("unexpected timestamp: %v", s.TakenAt())
	}
}

func TestGet_MissingKey(t *testing.T) {
	s := snapshot.Take(map[string]string{})
	_, ok := s.Get("MISSING")
	if ok {
		t.Fatal("expected miss for unknown key")
	}
}

func TestToMap_ReturnsIsolatedCopy(t *testing.T) {
	s := snapshot.Take(map[string]string{"A": "1"})
	m := s.ToMap()
	m["A"] = "mutated"
	v, _ := s.Get("A")
	if v != "1" {
		t.Fatal("snapshot was mutated through ToMap copy")
	}
}

func TestKeys_ReturnsAllKeys(t *testing.T) {
	s := snapshot.Take(map[string]string{"K1": "v1", "K2": "v2"})
	keys := s.Keys()
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
}
