// Package pipeline orchestrates the full secret-injection lifecycle:
// resolve → cache → inject → exec.
package pipeline

import (
	"context"
	"fmt"

	"github.com/your-org/vaultpipe/internal/environ"
	"github.com/your-org/vaultpipe/internal/redact"
	"github.com/your-org/vaultpipe/internal/resolver"
)

// SecretMap is a mapping from environment variable name to resolved secret value.
type SecretMap map[string]string

// Resolver fetches secrets by reference.
type Resolver interface {
	Resolve(ctx context.Context, ref string) (string, error)
}

// Runner executes a process with the given environment.
type Runner interface {
	Run(ctx context.Context, env []string) error
}

// Pipeline holds the dependencies for a single run.
type Pipeline struct {
	refs     []resolver.Ref
	res      Resolver
	runner   Runner
	redactor *redact.Redactor
}

// New constructs a Pipeline.
func New(refs []resolver.Ref, res Resolver, runner Runner, redactor *redact.Redactor) *Pipeline {
	return &Pipeline{
		refs:     refs,
		res:      res,
		runner:   runner,
		redactor: redactor,
	}
}

// Run resolves all secrets, builds the environment, then executes the child process.
func (p *Pipeline) Run(ctx context.Context, base *environ.Snapshot) error {
	secrets, err := p.resolveAll(ctx)
	if err != nil {
		return fmt.Errorf("pipeline: resolve: %w", err)
	}

	for _, v := range secrets {
		p.redactor.Add(v)
	}

	merged := base.Merge(secrets)

	if err := p.runner.Run(ctx, merged.Environ()); err != nil {
		return fmt.Errorf("pipeline: exec: %w", err)
	}
	return nil
}

func (p *Pipeline) resolveAll(ctx context.Context) (SecretMap, error) {
	out := make(SecretMap, len(p.refs))
	for _, ref := range p.refs {
		val, err := p.res.Resolve(ctx, ref.Path)
		if err != nil {
			return nil, fmt.Errorf("ref %q: %w", ref.Env, err)
		}
		out[ref.Env] = val
	}
	return out, nil
}
