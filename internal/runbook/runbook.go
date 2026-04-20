// Package runbook provides pre-flight validation checks before process
// execution, ensuring Vault connectivity and required secrets are accessible.
package runbook

import (
	"context"
	"fmt"
	"strings"

	"github.com/yourusername/vaultpipe/internal/config"
	"github.com/yourusername/vaultpipe/internal/vault"
)

// CheckResult holds the outcome of a single pre-flight check.
type CheckResult struct {
	Name    string
	Passed  bool
	Message string
}

// Runner executes a set of pre-flight checks.
type Runner struct {
	client *vault.Client
	cfg    *config.Config
}

// New returns a Runner for the given Vault client and config.
func New(client *vault.Client, cfg *config.Config) *Runner {
	return &Runner{client: client, cfg: cfg}
}

// Run executes all checks and returns their results.
// The returned error is non-nil only if a check cannot be evaluated at all.
func (r *Runner) Run(ctx context.Context) ([]CheckResult, error) {
	results := make([]CheckResult, 0, len(r.cfg.Secrets)+1)

	// Check 1: Vault reachability via token self-lookup.
	ok, msg := r.checkVaultHealth(ctx)
	results = append(results, CheckResult{
		Name:    "vault:reachable",
		Passed:  ok,
		Message: msg,
	})

	// Check 2: Each secret path is readable.
	for _, s := range r.cfg.Secrets {
		ok, msg := r.checkSecretReadable(ctx, s.Path)
		results = append(results, CheckResult{
			Name:    fmt.Sprintf("secret:%s", s.Path),
			Passed:  ok,
			Message: msg,
		})
	}

	return results, nil
}

// Summary returns a human-readable summary line and whether all checks passed.
func Summary(results []CheckResult) (string, bool) {
	var failed []string
	for _, r := range results {
		if !r.Passed {
			failed = append(failed, r.Name)
		}
	}
	if len(failed) == 0 {
		return fmt.Sprintf("all %d checks passed", len(results)), true
	}
	return fmt.Sprintf("%d/%d checks failed: %s", len(failed), len(results), strings.Join(failed, ", ")), false
}

func (r *Runner) checkVaultHealth(ctx context.Context) (bool, string) {
	if err := r.client.Ping(ctx); err != nil {
		return false, fmt.Sprintf("vault unreachable: %v", err)
	}
	return true, "vault is reachable"
}

func (r *Runner) checkSecretReadable(ctx context.Context, path string) (bool, string) {
	if _, err := r.client.Read(ctx, path); err != nil {
		return false, fmt.Sprintf("cannot read secret: %v", err)
	}
	return true, "secret is readable"
}
