package pipeline

import (
	"fmt"

	"github.com/your-org/vaultpipe/internal/audit"
	"github.com/your-org/vaultpipe/internal/cache"
	"github.com/your-org/vaultpipe/internal/config"
	"github.com/your-org/vaultpipe/internal/environ"
	"github.com/your-org/vaultpipe/internal/process"
	"github.com/your-org/vaultpipe/internal/redact"
	"github.com/your-org/vaultpipe/internal/resolver"
	"github.com/your-org/vaultpipe/internal/resolver/build"
	"github.com/your-org/vaultpipe/internal/vault"
)

// BuildOptions carries external dependencies that cannot be derived from
// the config file alone.
type BuildOptions struct {
	Args   []string
	Logger *audit.Logger
	Cache  *cache.Cache
}

// Build constructs a ready-to-run Pipeline from a loaded config.
func Build(cfg *config.Config, opts BuildOptions) (*Pipeline, *environ.Snapshot, error) {
	client, err := vault.NewClient(vault.Config{
		Address: cfg.Vault.Address,
		Token:   cfg.Vault.Token,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("pipeline build: vault client: %w", err)
	}

	refs := build.RefsFromConfig(cfg)

	res := resolver.New(client, opts.Cache, opts.Logger)

	runner, err := process.NewRunner(opts.Args)
	if err != nil {
		return nil, nil, fmt.Errorf("pipeline build: runner: %w", err)
	}

	redactor := redact.New()

	base := environ.Capture()

	p := New(refs, res, runner, redactor)
	return p, base, nil
}
