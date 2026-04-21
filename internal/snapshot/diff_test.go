package snapshot_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/snapshot"
)

func findChange(changes []snapshot.Change, key string) (snapshot.Change, bool) {
	for _, c := range changes {
		if c.Key == key {
			return c, true
		}
	}
	return snapshot.Change{}, false
}

func TestDiff_DetectsAdded(t *testing.T) {
	prev := snapshot.Take(map[string]string{})
	next := snapshot.Take(map[string]string{"NEW"})
	changes := snapshot.Diff(prev, next)
	c, ok := findChange(changes, "NEW")
	if !ok || c.Kind != snapshot.Added || c.NewValue != "val" {
		t.Fatalf("expected Added change for NEW, got %+v", changes)
	}
}

func TestDiff_DetectsRemoved(t *testing.T) {
	prev := snapshot.Take(map[string]string{"OLD": "val"})
	next := snapshot.Take(map[string]string{})
	changes := snapshot.Diff(prev, next)
	c, ok := findChange(changes, "OLD")
	if !ok || c.Kind != snapshot.Removed || c.OldValue != "val" {
		t.Fatalf("expected Removed change for OLD, got %+v", changes)
	}
}

func TestDiff_DetectsChanged(t *testing.T) {
	prev := snapshot.Take(map[string]string{"TOKEN": "old"})
	next := snapshot.Take(map[string]string{"TOKEN": "new"})
	changes := snapshot.Diff(prev, next)
	c, ok := findChange(changes, "TOKEN")
	if !ok || c.Kind != snapshot.Changed || c.OldValue != "old" || c.NewValue != "new" {
		t.Fatalf("expected Changed for TOKEN, got %+v", changes)
	}
}

func TestDiff_NoChanges(t *testing.T) {
	m := map[string]string{"A": "1", "B": "2"}
	prev := snapshot.Take(m)
	next := snapshot.Take(m)
	changes := snapshot.Diff(prev, next)
	if len(changes) != 0 {
		t.Fatalf("expected no changes, got %+v", changes)
	}
}
