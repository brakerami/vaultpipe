package vault

import (
	"testing"
)

func TestParseSecretPath_NoField(t *testing.T) {
	sp := ParseSecretPath("secret/data/myapp")
	if sp.Path != "secret/data/myapp" {
		t.Errorf("expected path %q, got %q", "secret/data/myapp", sp.Path)
	}
	if sp.Field != "" {
		t.Errorf("expected empty field, got %q", sp.Field)
	}
}

func TestParseSecretPath_WithField(t *testing.T) {
	sp := ParseSecretPath("secret/data/myapp#DB_PASSWORD")
	if sp.Path != "secret/data/myapp" {
		t.Errorf("expected path %q, got %q", "secret/data/myapp", sp.Path)
	}
	if sp.Field != "DB_PASSWORD" {
		t.Errorf("expected field %q, got %q", "DB_PASSWORD", sp.Field)
	}
}

func TestParseSecretPath_MultipleHashes(t *testing.T) {
	// Only the first '#' should split path from field
	sp := ParseSecretPath("secret/data/myapp#field#extra")
	if sp.Path != "secret/data/myapp" {
		t.Errorf("expected path %q, got %q", "secret/data/myapp", sp.Path)
	}
	if sp.Field != "field#extra" {
		t.Errorf("expected field %q, got %q", "field#extra", sp.Field)
	}
}

func TestParseSecretPath_EmptyString(t *testing.T) {
	sp := ParseSecretPath("")
	if sp.Path != "" {
		t.Errorf("expected empty path, got %q", sp.Path)
	}
	if sp.Field != "" {
		t.Errorf("expected empty field, got %q", sp.Field)
	}
}
