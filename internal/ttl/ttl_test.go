package ttl_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/ttl"
)

func TestParse_EmptyString_ReturnsDefault(t *testing.T) {
	d, err := ttl.Parse("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d != ttl.DefaultTTL {
		t.Errorf("expected %v, got %v", ttl.DefaultTTL, d)
	}
}

func TestParse_ValidDuration(t *testing.T) {
	cases := []struct {
		input    string
		expected time.Duration
	}{
		{"10s", 10 * time.Second},
		{"2m", 2 * time.Minute},
		{"1h", 1 * time.Hour},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			d, err := ttl.Parse(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if d != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, d)
			}
		})
	}
}

func TestParse_InvalidString(t *testing.T) {
	_, err := ttl.Parse("not-a-duration")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParse_TooShort(t *testing.T) {
	_, err := ttl.Parse("500ms")
	if err != ttl.ErrTTLTooShort {
		t.Errorf("expected ErrTTLTooShort, got %v", err)
	}
}

func TestParse_TooLong(t *testing.T) {
	_, err := ttl.Parse("25h")
	if err != ttl.ErrTTLTooLong {
		t.Errorf("expected ErrTTLTooLong, got %v", err)
	}
}

func TestValidate_BoundaryValues(t *testing.T) {
	if _, err := ttl.Validate(ttl.MinTTL); err != nil {
		t.Errorf("MinTTL should be valid: %v", err)
	}
	if _, err := ttl.Validate(ttl.MaxTTL); err != nil {
		t.Errorf("MaxTTL should be valid: %v", err)
	}
}

func TestMustParse_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic, got none")
		}
	}()
	ttl.MustParse("bad")
}
