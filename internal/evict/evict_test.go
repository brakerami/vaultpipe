package evict

import (
	"testing"
)

func TestNew_PanicsOnZeroCap(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for zero capacity")
		}
	}()
	New(0)
}

func TestSet_And_Get(t *testing.T) {
	l := New(3)
	l.Set("a", "1")
	l.Set("b", "2")

	v, ok := l.Get("a")
	if !ok || v != "1" {
		t.Fatalf("expected (1, true), got (%q, %v)", v, ok)
	}
}

func TestGet_MissingKey(t *testing.T) {
	l := New(2)
	_, ok := l.Get("missing")
	if ok {
		t.Fatal("expected false for missing key")
	}
}

func TestSet_UpdateExisting(t *testing.T) {
	l := New(2)
	l.Set("k", "old")
	l.Set("k", "new")

	v, _ := l.Get("k")
	if v != "new" {
		t.Fatalf("expected 'new', got %q", v)
	}
	if l.Len() != 1 {
		t.Fatalf("expected len 1, got %d", l.Len())
	}
}

func TestEviction_RemovesLRU(t *testing.T) {
	evicted := map[string]string{}
	l := New(2)
	l.Evicted = func(k, v string) { evicted[k] = v }

	l.Set("a", "1")
	l.Set("b", "2")
	// access "a" so "b" becomes LRU
	l.Get("a")
	l.Set("c", "3") // should evict "b"

	if _, ok := evicted["b"]; !ok {
		t.Fatal("expected 'b' to be evicted")
	}
	if _, ok := l.items["b"]; ok {
		t.Fatal("'b' should not be in cache after eviction")
	}
}

func TestRemove_DeletesKey(t *testing.T) {
	l := New(3)
	l.Set("x", "val")
	l.Remove("x")

	_, ok := l.Get("x")
	if ok {
		t.Fatal("expected key to be removed")
	}
}

func TestRemove_NoopOnMissing(t *testing.T) {
	l := New(2)
	l.Remove("nonexistent") // must not panic
}

func TestLen_TracksSize(t *testing.T) {
	l := New(5)
	if l.Len() != 0 {
		t.Fatalf("expected 0, got %d", l.Len())
	}
	l.Set("a", "1")
	l.Set("b", "2")
	if l.Len() != 2 {
		t.Fatalf("expected 2, got %d", l.Len())
	}
	l.Remove("a")
	if l.Len() != 1 {
		t.Fatalf("expected 1, got %d", l.Len())
	}
}
