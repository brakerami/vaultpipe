package window

import (
	"testing"
	"time"
)

func TestNew_DefaultsBuckets(t *testing.T) {
	w := New(time.Second, 0)
	if len(w.counts) != 1 {
		t.Fatalf("expected 1 bucket, got %d", len(w.counts))
	}
}

func TestAdd_IncreasesCount(t *testing.T) {
	w := New(time.Second, 10)
	w.Add(3)
	w.Add(2)
	if got := w.Count(); got != 5 {
		t.Fatalf("expected 5, got %d", got)
	}
}

func TestCount_ExcludesExpiredBuckets(t *testing.T) {
	now := time.Now()
	w := New(100*time.Millisecond, 2)
	// inject fake clock
	w.now = func() time.Time { return now }
	w.Add(10)

	// advance time beyond the window
	w.now = func() time.Time { return now.Add(200 * time.Millisecond) }
	if got := w.Count(); got != 0 {
		t.Fatalf("expected 0 after expiry, got %d", got)
	}
}

func TestCount_IncludesRecentBuckets(t *testing.T) {
	now := time.Now()
	w := New(1*time.Second, 4)
	w.now = func() time.Time { return now }
	w.Add(7)
	w.now = func() time.Time { return now.Add(100 * time.Millisecond) }
	w.Add(3)
	if got := w.Count(); got != 10 {
		t.Fatalf("expected 10, got %d", got)
	}
}

func TestReset_ClearsAll(t *testing.T) {
	w := New(time.Second, 4)
	w.Add(99)
	w.Reset()
	if got := w.Count(); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestCount_ZeroWhenEmpty(t *testing.T) {
	w := New(time.Second, 4)
	if got := w.Count(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}
