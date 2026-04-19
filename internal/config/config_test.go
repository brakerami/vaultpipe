package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writeTemp: %v", err)
	}
	return p
}

func TestLoadFile_Valid(t *testing.T) {
	p := writeTemp(t, `secrets:
  - env: DB_PASSWORD
    path: secret/data/db#password
  - env: API_KEY
    path: secret/data/api#key
`)
	cfg, err := LoadFile(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Secrets) != 2 {
		t.Fatalf("expected 2 secrets, got %d", len(cfg.Secrets))
	}
	if cfg.Secrets[0].EnvVar != "DB_PASSWORD" {
		t.Errorf("expected DB_PASSWORD, got %q", cfg.Secrets[0].EnvVar)
	}
}

func TestLoadFile_MissingFile(t *testing.T) {
	_, err := LoadFile("/nonexistent/path/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadFile_MissingEnv(t *testing.T) {
	p := writeTemp(t, `secrets:
  - env: ""
    path: secret/data/db#password
`)
	_, err := LoadFile(p)
	if err == nil {
		t.Fatal("expected validation error for empty env")
	}
}

func TestLoadFile_MissingPath(t *testing.T) {
	p := writeTemp(t, `secrets:
  - env: DB_PASSWORD
    path: ""
`)
	_, err := LoadFile(p)
	if err == nil {
		t.Fatal("expected validation error for empty path")
	}
}

func TestLoadFile_EmptySecrets(t *testing.T) {
	p := writeTemp(t, `secrets: []\n`)
	cfg, err := LoadFile(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Secrets) != 0 {
		t.Errorf("expected 0 secrets, got %d", len(cfg.Secrets))
	}
}
