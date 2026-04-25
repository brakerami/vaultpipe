package truncate_test

import (
	"strings"
	"testing"

	"github.com/yourusername/vaultpipe/internal/truncate"
)

func TestNew_DefaultsWhenZero(t *testing.T) {
	tr := truncate.New(0)
	if tr == nil {
		t.Fatal("expected non-nil Truncator")
	}
	// A value well within the default limit should pass through unchanged.
	out, err := tr.Apply("KEY", "hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "hello" {
		t.Fatalf("expected %q, got %q", "hello", out)
	}
}

func TestApply_WithinLimit_NoError(t *testing.T) {
	tr := truncate.New(10)
	out, err := tr.Apply("K", "hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "hello" {
		t.Fatalf("expected %q, got %q", "hello", out)
	}
}

func TestApply_ExceedsLimit_Truncates(t *testing.T) {
	tr := truncate.New(5)
	out, err := tr.Apply("SECRET", "hello world")
	if err == nil {
		t.Fatal("expected ErrTruncated, got nil")
	}
	te, ok := err.(*truncate.ErrTruncated)
	if !ok {
		t.Fatalf("expected *ErrTruncated, got %T", err)
	}
	if te.Key != "SECRET" {
		t.Errorf("expected key %q, got %q", "SECRET", te.Key)
	}
	if len(out) > 5 {
		t.Errorf("output length %d exceeds limit 5", len(out))
	}
	if out != "hello" {
		t.Errorf("expected %q, got %q", "hello", out)
	}
}

func TestApply_PreservesUTF8Boundary(t *testing.T) {
	// "é" is 2 bytes (0xC3 0xA9). With limit=3 and value="aéb" (4 bytes),
	// the truncator must not split the rune.
	tr := truncate.New(3)
	out, err := tr.Apply("K", "aéb")
	if err == nil {
		t.Fatal("expected truncation error")
	}
	if !strings.HasPrefix("aéb", out) {
		t.Errorf("truncated value %q is not a prefix of original", out)
	}
	if !isValidUTF8(out) {
		t.Errorf("truncated value %q is not valid UTF-8", out)
	}
}

func TestMap_TruncatesOversizedEntries(t *testing.T) {
	tr := truncate.New(4)
	env := map[string]string{
		"SHORT": "hi",
		"LONG":  "toolongvalue",
	}
	out, errs := tr.Map(env)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
	if out["SHORT"] != "hi" {
		t.Errorf("SHORT changed unexpectedly: %q", out["SHORT"])
	}
	if len(out["LONG"]) > 4 {
		t.Errorf("LONG not truncated: len=%d", len(out["LONG"]))
	}
}

func TestMap_AllWithinLimit_NoErrors(t *testing.T) {
	tr := truncate.New(100)
	env := map[string]string{"A": "foo", "B": "bar"}
	_, errs := tr.Map(env)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func isValidUTF8(s string) bool {
	return strings.ToValidUTF8(s, "\uFFFD") == s
}
