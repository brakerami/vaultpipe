package buffer_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/buffer"
)

func TestNew_PanicsOnZeroCap(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for zero capacity")
		}
	}()
	buffer.New(0)
}

func TestAdd_And_Len(t *testing.T) {
	r := buffer.New(3)
	if r.Len() != 0 {
		t.Fatalf("expected 0, got %d", r.Len())
	}
	r.Add("a")
	r.Add("b")
	if r.Len() != 2 {
		t.Fatalf("expected 2, got %d", r.Len())
	}
}

func TestSnapshot_ChronologicalOrder(t *testing.T) {
	r := buffer.New(4)
	msgs := []string{"first", "second", "third"}
	for _, m := range msgs {
		r.Add(m)
	}
	snap := r.Snapshot()
	if len(snap) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(snap))
	}
	for i, e := range snap {
		if e.Message != msgs[i] {
			t.Errorf("entry %d: want %q, got %q", i, msgs[i], e.Message)
		}
	}
}

func TestRing_OverwritesOldest(t *testing.T) {
	r := buffer.New(2)
	r.Add("old")
	r.Add("keep")
	r.Add("new")
	snap := r.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(snap))
	}
	if snap[0].Message != "keep" {
		t.Errorf("expected 'keep', got %q", snap[0].Message)
	}
	if snap[1].Message != "new" {
		t.Errorf("expected 'new', got %q", snap[1].Message)
	}
}

func TestSnapshot_SetsTimestamp(t *testing.T) {
	before := time.Now()
	r := buffer.New(2)
	r.Add("ts-check")
	after := time.Now()
	snap := r.Snapshot()
	if snap[0].At.Before(before) || snap[0].At.After(after) {
		t.Error("timestamp out of expected range")
	}
}

func TestReset_ClearsBuffer(t *testing.T) {
	r := buffer.New(4)
	r.Add("x")
	r.Add("y")
	r.Reset()
	if r.Len() != 0 {
		t.Fatalf("expected 0 after reset, got %d", r.Len())
	}
	if snap := r.Snapshot(); snap != nil {
		t.Errorf("expected nil snapshot after reset, got %v", snap)
	}
}

func TestSnapshot_EmptyReturnsNil(t *testing.T) {
	r := buffer.New(5)
	if snap := r.Snapshot(); snap != nil {
		t.Errorf("expected nil, got %v", snap)
	}
}
