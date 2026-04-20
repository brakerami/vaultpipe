package environ_test

import (
	"testing"

	"github.com/your-org/vaultpipe/internal/environ"
)

func TestFromMap_StoresValues(t *testing.T) {
	s := environ.FromMap(map[string]string{
		"FOO": "bar",
		"BAZ": "qux",
	})

	v, ok := s.Get("FOO")
	if !ok || v != "bar" {
		t.Errorf("expected FOO=bar, got %q (ok=%v)", v, ok)
	}
}

func TestFromMap_IsolatesOriginal(t *testing.T) {
	orig := map[string]string{"KEY": "original"}
	s := environ.FromMap(orig)
	orig["KEY"] = "mutated"

	v, _ := s.Get("KEY")
	if v != "original" {
		t.Errorf("snapshot should be isolated from source map, got %q", v)
	}
}

func TestGet_MissingKey(t *testing.T) {
	s := environ.FromMap(map[string]string{})
	_, ok := s.Get("MISSING")
	if ok {
		t.Error("expected ok=false for missing key")
	}
}

func TestEnviron_FormatsPairs(t *testing.T) {
	s := environ.FromMap(map[string]string{"A": "1"})
	env := s.Environ()
	if len(env) != 1 || env[0] != "A=1" {
		t.Errorf("unexpected Environ output: %v", env)
	}
}

func TestMerge_AppliesOverrides(t *testing.T) {
	base := environ.FromMap(map[string]string{
		"BASE": "keep",
		"OVER": "old",
	})
	merged := base.Merge(map[string]string{"OVER": "new", "EXTRA": "added"})

	if v, _ := merged.Get("BASE"); v != "keep" {
		t.Errorf("BASE should be kept, got %q", v)
	}
	if v, _ := merged.Get("OVER"); v != "new" {
		t.Errorf("OVER should be overridden, got %q", v)
	}
	if v, _ := merged.Get("EXTRA"); v != "added" {
		t.Errorf("EXTRA should be added, got %q", v)
	}
}

func TestMerge_DoesNotMutateBase(t *testing.T) {
	base := environ.FromMap(map[string]string{"KEY": "original"})
	base.Merge(map[string]string{"KEY": "changed"})

	if v, _ := base.Get("KEY"); v != "original" {
		t.Errorf("base snapshot should not be mutated, got %q", v)
	}
}

func TestLen_ReturnsCount(t *testing.T) {
	s := environ.FromMap(map[string]string{"A": "1", "B": "2", "C": "3"})
	if s.Len() != 3 {
		t.Errorf("expected Len()=3, got %d", s.Len())
	}
}

func TestKeys_ReturnsAllKeys(t *testing.T) {
	s := environ.FromMap(map[string]string{"X": "1", "Y": "2"})
	keys := s.Keys()
	if len(keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(keys))
	}
}
