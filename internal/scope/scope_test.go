package scope_test

import (
	"testing"

	"github.com/your-org/vaultpipe/internal/scope"
)

func TestNew_EmptyPrefix_ReturnsError(t *testing.T) {
	_, err := scope.New("", "_")
	if err == nil {
		t.Fatal("expected error for empty prefix")
	}
}

func TestNew_InvalidPrefix_ReturnsError(t *testing.T) {
	for _, bad := range []string{"bad prefix", "bad=prefix", "bad\nprefix"} {
		_, err := scope.New(bad, "_")
		if err == nil {
			t.Fatalf("expected error for prefix %q", bad)
		}
	}
}

func TestNew_DefaultSeparator(t *testing.T) {
	s, err := scope.New("APP", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := s.Apply("TOKEN")
	want := "APP_TOKEN"
	if got != want {
		t.Fatalf("Apply() = %q, want %q", got, want)
	}
}

func TestApply_UppercasesKey(t *testing.T) {
	s, _ := scope.New("svc", "_")
	got := s.Apply("db_password")
	want := "SVC_DB_PASSWORD"
	if got != want {
		t.Fatalf("Apply() = %q, want %q", got, want)
	}
}

func TestApply_EmptyKey_ReturnsPrefix(t *testing.T) {
	s, _ := scope.New("SVC", "_")
	got := s.Apply("")
	if got != "SVC" {
		t.Fatalf("Apply(\"\") = %q, want %q", got, "SVC")
	}
}

func TestStrip_MatchingKey_ReturnsStripped(t *testing.T) {
	s, _ := scope.New("APP", "_")
	got, ok := s.Strip("APP_SECRET")
	if !ok {
		t.Fatal("expected match")
	}
	if got != "SECRET" {
		t.Fatalf("Strip() = %q, want %q", got, "SECRET")
	}
}

func TestStrip_NonMatchingKey_ReturnsFalse(t *testing.T) {
	s, _ := scope.New("APP", "_")
	got, ok := s.Strip("OTHER_SECRET")
	if ok {
		t.Fatal("expected no match")
	}
	if got != "OTHER_SECRET" {
		t.Fatalf("Strip() = %q, want original", got)
	}
}

func TestMap_AppliesPrefixToAllKeys(t *testing.T) {
	s, _ := scope.New("NS", "_")
	src := map[string]string{"key": "val1", "other": "val2"}
	out := s.Map(src)
	for _, k := range []string{"NS_KEY", "NS_OTHER"} {
		if _, ok := out[k]; !ok {
			t.Fatalf("expected key %q in output", k)
		}
	}
	if len(out) != len(src) {
		t.Fatalf("Map() len = %d, want %d", len(out), len(src))
	}
}

func TestPrefix_ReturnsUppercasedPrefix(t *testing.T) {
	s, _ := scope.New("myapp", "_")
	if s.Prefix() != "MYAPP" {
		t.Fatalf("Prefix() = %q, want %q", s.Prefix(), "MYAPP")
	}
}
