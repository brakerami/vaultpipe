package pipeline_test

import (
	"testing"

	"github.com/your-org/vaultpipe/internal/audit"
	"github.com/your-org/vaultpipe/internal/cache"
	"github.com/your-org/vaultpipe/internal/config"
	"github.com/your-org/vaultpipe/internal/pipeline"
)

func minimalConfig(addr, token string) *config.Config {
	return &config.Config{
		Vault: config.VaultConfig{
			Address: addr,
			Token:   token,
		},
		Secrets: []config.SecretEntry{
			{Env: "FOO", Path: "secret/foo"},
		},
	}
}

func TestBuild_MissingVaultAddress(t *testing.T) {
	cfg := minimalConfig("", "tok")
	opts := pipeline.BuildOptions{
		Args:   []string{"echo", "hi"},
		Logger: audit.NewLogger(nil),
		Cache:  cache.New(0),
	}
	_, _, err := pipeline.Build(cfg, opts)
	if err == nil {
		t.Fatal("expected error for missing vault address")
	}
}

func TestBuild_MissingVaultToken(t *testing.T) {
	cfg := minimalConfig("http://127.0.0.1:8200", "")
	opts := pipeline.BuildOptions{
		Args:   []string{"echo"},
		Logger: audit.NewLogger(nil),
		Cache:  cache.New(0),
	}
	_, _, err := pipeline.Build(cfg, opts)
	if err == nil {
		t.Fatal("expected error for missing vault token")
	}
}

func TestBuild_NoArgs(t *testing.T) {
	cfg := minimalConfig("http://127.0.0.1:8200", "root")
	opts := pipeline.BuildOptions{
		Args:   []string{},
		Logger: audit.NewLogger(nil),
		Cache:  cache.New(0),
	}
	_, _, err := pipeline.Build(cfg, opts)
	if err == nil {
		t.Fatal("expected error for empty args")
	}
}

func TestBuild_ValidConfig(t *testing.T) {
	cfg := minimalConfig("http://127.0.0.1:8200", "root")
	opts := pipeline.BuildOptions{
		Args:   []string{"env"},
		Logger: audit.NewLogger(nil),
		Cache:  cache.New(0),
	}
	p, base, err := pipeline.Build(cfg, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil pipeline")
	}
	if base == nil {
		t.Fatal("expected non-nil base snapshot")
	}
}
