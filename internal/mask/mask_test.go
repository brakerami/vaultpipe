package mask_test

import (
	"testing"

	"github.com/yourorg/vaultpipe/internal/mask"
)

func TestMask_RedactsSecret(t *testing.T) {
	m := mask.New()
	m.Add("s3cr3t")

	got := m.Mask("password is s3cr3t, keep it safe")
	want := "password is [REDACTED], keep it safe"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMask_MultipleSecrets(t *testing.T) {
	m := mask.New()
	m.Add("alpha")
	m.Add("beta")

	got := m.Mask("alpha and beta are secrets")
	if got == "alpha and beta are secrets" {
		t.Error("expected secrets to be redacted")
	}
}

func TestMask_EmptyStringIgnored(t *testing.T) {
	m := mask.New()
	m.Add("")
	if m.Len() != 0 {
		t.Errorf("expected 0 secrets, got %d", m.Len())
	}
}

func TestMask_NoSecretsRegistered(t *testing.T) {
	m := mask.New()
	input := "nothing to hide"
	got := m.Mask(input)
	if got != input {
		t.Errorf("got %q, want %q", got, input)
	}
}

func TestMask_Len(t *testing.T) {
	m := mask.New()
	m.Add("one")
	m.Add("two")
	m.Add("one") // duplicate
	if m.Len() != 2 {
		t.Errorf("expected 2, got %d", m.Len())
	}
}
