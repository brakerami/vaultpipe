package tag_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/tag"
)

func TestAdd_And_Get(t *testing.T) {
	s := tag.New()
	if err := s.Add("env", "prod"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, ok := s.Get("env")
	if !ok || v != "prod" {
		t.Fatalf("expected env=prod, got %q ok=%v", v, ok)
	}
}

func TestAdd_EmptyKey_ReturnsError(t *testing.T) {
	s := tag.New()
	if err := s.Add("", "val"); err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestAdd_EmptyValue_ReturnsError(t *testing.T) {
	s := tag.New()
	if err := s.Add("env", ""); err == nil {
		t.Fatal("expected error for empty value")
	}
}

func TestAdd_ReservedCharInKey_ReturnsError(t *testing.T) {
	s := tag.New()
	for _, bad := range []string{"a=b", "a,b"} {
		if err := s.Add(bad, "v"); err == nil {
			t.Fatalf("expected error for key %q", bad)
		}
	}
}

func TestGet_MissingKey_ReturnsFalse(t *testing.T) {
	s := tag.New()
	_, ok := s.Get("missing")
	if ok {
		t.Fatal("expected ok=false for missing key")
	}
}

func TestHas_ReportsPresence(t *testing.T) {
	s := tag.New()
	_ = s.Add("team", "platform")
	if !s.Has("team") {
		t.Fatal("expected Has to return true")
	}
	if s.Has("other") {
		t.Fatal("expected Has to return false for absent key")
	}
}

func TestLen_TracksCount(t *testing.T) {
	s := tag.New()
	if s.Len() != 0 {
		t.Fatalf("expected 0, got %d", s.Len())
	}
	_ = s.Add("a", "1")
	_ = s.Add("b", "2")
	if s.Len() != 2 {
		t.Fatalf("expected 2, got %d", s.Len())
	}
}

func TestString_SerialisesPairs(t *testing.T) {
	s := tag.New()
	_ = s.Add("env", "prod")
	_ = s.Add("region", "us-east")
	got := s.String()
	want := "env=prod,region=us-east"
	if got != want {
		t.Fatalf("String() = %q; want %q", got, want)
	}
}

func TestParse_RoundTrip(t *testing.T) {
	original := "env=prod,region=us-east"
	s, err := tag.Parse(original)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if s.String() != original {
		t.Fatalf("round-trip mismatch: %q", s.String())
	}
}

func TestParse_EmptyString_ReturnsEmptySet(t *testing.T) {
	s, err := tag.Parse("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Len() != 0 {
		t.Fatalf("expected empty set, got len=%d", s.Len())
	}
}

func TestParse_MalformedPair_ReturnsError(t *testing.T) {
	if _, err := tag.Parse("noequalssign"); err == nil {
		t.Fatal("expected error for malformed pair")
	}
}
