package redact

import (
	"errors"
	"testing"
)

func TestScrub_ReplacesSecret(t *testing.T) {
	r := New([]string{"s3cr3t"})
	got := r.Scrub("the value is s3cr3t ok")
	want := "the value is *** ok"
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestScrub_MultipleSecrets(t *testing.T) {
	r := New([]string{"alpha", "beta"})
	got := r.Scrub("alpha and beta")
	if got != "*** and ***" {
		t.Errorf("unexpected: %q", got)
	}
}

func TestScrub_EmptySecretIgnored(t *testing.T) {
	r := New([]string{"", "real"})
	if len(r.secrets) != 1 {
		t.Fatalf("expected 1 secret, got %d", len(r.secrets))
	}
}

func TestScrub_NoSecrets(t *testing.T) {
	r := New(nil)
	got := r.Scrub("plain text")
	if got != "plain text" {
		t.Errorf("unexpected: %q", got)
	}
}

func TestAdd_RegistersSecret(t *testing.T) {
	r := New(nil)
	r.Add("newval")
	got := r.Scrub("contains newval here")
	if got != "contains *** here" {
		t.Errorf("unexpected: %q", got)
	}
}

func TestScrubError_NilReturnsNil(t *testing.T) {
	r := New([]string{"x"})
	if r.ScrubError(nil) != nil {
		t.Error("expected nil")
	}
}

func TestScrubError_RedactsMessage(t *testing.T) {
	r := New([]string{"topsecret"})
	err := errors.New("failed with topsecret value")
	got := r.ScrubError(err).Error()
	if got != "failed with *** value" {
		t.Errorf("unexpected: %q", got)
	}
}
