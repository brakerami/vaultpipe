package clamp_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultpipe/internal/clamp"
)

func TestInt_WithinRange(t *testing.T) {
	out, err := clamp.Int(5, 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != 5 {
		t.Fatalf("expected 5, got %d", out)
	}
}

func TestInt_BelowMin(t *testing.T) {
	out, err := clamp.Int(-3, 0, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != 0 {
		t.Fatalf("expected 0, got %d", out)
	}
}

func TestInt_AboveMax(t *testing.T) {
	out, err := clamp.Int(99, 0, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != 10 {
		t.Fatalf("expected 10, got %d", out)
	}
}

func TestInt_InvalidBounds(t *testing.T) {
	_, err := clamp.Int(5, 10, 1)
	if err == nil {
		t.Fatal("expected error for invalid bounds")
	}
}

func TestMustInt_PanicsOnInvalidBounds(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	clamp.MustInt(5, 10, 1)
}

func TestDuration_WithinRange(t *testing.T) {
	out, err := clamp.Duration(5*time.Second, 1*time.Second, 10*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != 5*time.Second {
		t.Fatalf("expected 5s, got %s", out)
	}
}

func TestDuration_BelowMin(t *testing.T) {
	out, err := clamp.Duration(100*time.Millisecond, 1*time.Second, 10*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != 1*time.Second {
		t.Fatalf("expected 1s, got %s", out)
	}
}

func TestDuration_AboveMax(t *testing.T) {
	out, err := clamp.Duration(1*time.Hour, 1*time.Second, 10*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != 10*time.Second {
		t.Fatalf("expected 10s, got %s", out)
	}
}

func TestDuration_InvalidBounds(t *testing.T) {
	_, err := clamp.Duration(5*time.Second, 10*time.Second, 1*time.Second)
	if err == nil {
		t.Fatal("expected error for invalid bounds")
	}
}

func TestMustDuration_PanicsOnInvalidBounds(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	clamp.MustDuration(5*time.Second, 10*time.Second, 1*time.Second)
}
