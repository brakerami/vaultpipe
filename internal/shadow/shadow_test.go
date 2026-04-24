package shadow_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/shadow"
)

func TestSet_And_Get(t *testing.T) {
	s := shadow.New()
	s.Set("DB_PASS", "secret123")

	e, ok := s.Get("DB_PASS")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Value != "secret123" {
		t.Errorf("got %q, want %q", e.Value, "secret123")
	}
	if e.RecordedAt.IsZero() {
		t.Error("RecordedAt should not be zero")
	}
}

func TestGet_MissingKey(t *testing.T) {
	s := shadow.New()
	_, ok := s.Get("MISSING")
	if ok {
		t.Error("expected missing key to return false")
	}
}

func TestDrifted_NoShadow_ReturnsFalse(t *testing.T) {
	s := shadow.New()
	if s.Drifted("KEY", "value") {
		t.Error("expected false when no shadow exists")
	}
}

func TestDrifted_SameValue_ReturnsFalse(t *testing.T) {
	s := shadow.New()
	s.Set("KEY", "stable")
	if s.Drifted("KEY", "stable") {
		t.Error("expected false when values match")
	}
}

func TestDrifted_ChangedValue_ReturnsTrue(t *testing.T) {
	s := shadow.New()
	s.Set("KEY", "old")
	if !s.Drifted("KEY", "new") {
		t.Error("expected true when value has changed")
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	s := shadow.New()
	s.Set("KEY", "val")
	s.Delete("KEY")
	_, ok := s.Get("KEY")
	if ok {
		t.Error("expected entry to be removed after Delete")
	}
}

func TestKeys_ReturnsAllTracked(t *testing.T) {
	s := shadow.New()
	s.Set("A", "1")
	s.Set("B", "2")
	s.Set("C", "3")

	keys := s.Keys()
	if len(keys) != 3 {
		t.Errorf("expected 3 keys, got %d", len(keys))
	}
}
