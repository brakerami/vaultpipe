package sanitize_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/sanitize"
)

func TestValue_PlainString_Unchanged(t *testing.T) {
	input := "hunter2"
	got := sanitize.Value(input)
	if got != input {
		t.Fatalf("expected %q, got %q", input, got)
	}
}

func TestValue_StripsNullByte(t *testing.T) {
	input := "secret\x00value"
	want := "secretvalue"
	got := sanitize.Value(input)
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestValue_StripsEscapeSequence(t *testing.T) {
	// ESC [ 1 m is a common ANSI bold sequence
	input := "\x1b[1mbold\x1b[0m"
	want := "[1mbold[0m"
	got := sanitize.Value(input)
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestValue_PreservesTabAndSpace(t *testing.T) {
	input := "hello world\there"
	got := sanitize.Value(input)
	if got != input {
		t.Fatalf("expected %q, got %q", input, got)
	}
}

func TestValue_StripsNewline(t *testing.T) {
	input := "line1\nline2"
	want := "line1line2"
	got := sanitize.Value(input)
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestValue_EmptyString(t *testing.T) {
	got := sanitize.Value("")
	if got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestMap_SanitizesAllValues(t *testing.T) {
	input := map[string]string{
		"KEY_A": "clean",
		"KEY_B": "dirty\x00value",
		"KEY_C": "also\ndirty",
	}
	want := map[string]string{
		"KEY_A": "clean",
		"KEY_B": "dirtyvalue",
		"KEY_C": "alsodirty",
	}
	got := sanitize.Map(input)
	for k, wv := range want {
		if got[k] != wv {
			t.Errorf("key %s: expected %q, got %q", k, wv, got[k])
		}
	}
}

func TestMap_DoesNotMutateOriginal(t *testing.T) {
	orig := map[string]string{"K": "v\x00alue"}
	_ = sanitize.Map(orig)
	if orig["K"] != "v\x00alue" {
		t.Fatal("original map was mutated")
	}
}
