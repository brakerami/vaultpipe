package expire_test

import (
	"sync"
	"testing"
	"time"

	"vaultpipe/internal/expire"
)

func TestNew_PanicsOnZeroThreshold(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for zero threshold")
		}
	}()
	expire.New(0, func(string, time.Duration) {})
}

func TestNew_PanicsOnNilHandler(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil handler")
		}
	}()
	expire.New(time.Minute, nil)
}

func TestCheck_FiresHandlerWhenBelowThreshold(t *testing.T) {
	var mu sync.Mutex
	fired := map[string]time.Duration{}

	tr := expire.New(5*time.Minute, func(key string, rem time.Duration) {
		mu.Lock()
		fired[key] = rem
		mu.Unlock()
	})

	now := time.Now()
	tr.Track("soon", now.Add(2*time.Minute))   // within threshold
	tr.Track("later", now.Add(10*time.Minute)) // outside threshold

	tr.Check(now)

	mu.Lock()
	defer mu.Unlock()
	if _, ok := fired["soon"]; !ok {
		t.Error("expected handler to fire for 'soon'")
	}
	if _, ok := fired["later"]; ok {
		t.Error("handler should not fire for 'later'")
	}
}

func TestCheck_RemovesExpiredEntry(t *testing.T) {
	tr := expire.New(time.Minute, func(string, time.Duration) {})

	now := time.Now()
	tr.Track("gone", now.Add(-1*time.Second)) // already expired

	if tr.Len() != 1 {
		t.Fatalf("expected 1 entry before check, got %d", tr.Len())
	}

	tr.Check(now)

	if tr.Len() != 0 {
		t.Errorf("expected expired entry to be removed, len=%d", tr.Len())
	}
}

func TestRemove_StopsTracking(t *testing.T) {
	called := false
	tr := expire.New(time.Hour, func(string, time.Duration) { called = true })

	now := time.Now()
	tr.Track("key", now.Add(30*time.Minute))
	tr.Remove("key")
	tr.Check(now.Add(29 * time.Minute))

	if called {
		t.Error("handler should not fire after Remove")
	}
}

func TestTrack_ReplacesExistingEntry(t *testing.T) {
	var mu sync.Mutex
	count := 0
	tr := expire.New(5*time.Minute, func(string, time.Duration) {
		mu.Lock()
		count++
		mu.Unlock()
	})

	now := time.Now()
	tr.Track("k", now.Add(1*time.Minute)) // within threshold
	tr.Track("k", now.Add(1*time.Minute)) // replace — still one entry

	if tr.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", tr.Len())
	}

	tr.Check(now)

	mu.Lock()
	defer mu.Unlock()
	if count != 1 {
		t.Errorf("expected handler called once, got %d", count)
	}
}
