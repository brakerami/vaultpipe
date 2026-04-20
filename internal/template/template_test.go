package template_test

import (
	"errors"
	"testing"

	"github.com/your-org/vaultpipe/internal/template"
)

// staticResolver resolves from a fixed map.
type staticResolver map[string]string

func (s staticResolver) Resolve(path, field string) (string, error) {
	key := path + "#" + field
	if v, ok := s[key]; ok {
		return v, nil
	}
	return "", errors.New("not found: " + key)
}

func TestExpand_NoPlaceholders(t *testing.T) {
	out, err := template.Expand("hello world", staticResolver{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "hello world" {
		t.Fatalf("expected unchanged string, got %q", out)
	}
}

func TestExpand_SingleRef(t *testing.T) {
	r := staticResolver{"secret/db#password": "s3cr3t"}
	out, err := template.Expand("{{vault:secret/db#password}}", r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "s3cr3t" {
		t.Fatalf("expected s3cr3t, got %q", out)
	}
}

func TestExpand_MultipleRefs(t *testing.T) {
	r := staticResolver{
		"secret/db#user":     "admin",
		"secret/db#password": "pass",
	}
	out, err := template.Expand("{{vault:secret/db#user}}:{{vault:secret/db#password}}", r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "admin:pass" {
		t.Fatalf("expected admin:pass, got %q", out)
	}
}

func TestExpand_ResolveError(t *testing.T) {
	_, err := template.Expand("{{vault:secret/missing#key}}", staticResolver{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestExpand_MissingField(t *testing.T) {
	_, err := template.Expand("{{vault:secret/db}}", staticResolver{})
	if err == nil {
		t.Fatal("expected error for missing field separator")
	}
}

func TestHasRefs(t *testing.T) {
	if template.HasRefs("no refs here") {
		t.Fatal("expected false")
	}
	if !template.HasRefs("prefix-{{vault:a/b#c}}-suffix") {
		t.Fatal("expected true")
	}
}
