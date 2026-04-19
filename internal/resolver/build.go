package resolver

import (
	"fmt"

	"github.com/your-org/vaultpipe/internal/config"
	"github.com/your-org/vaultpipe/internal/vault"
)

// RefsFromConfig converts a loaded config into a slice of SecretRefs.
func RefsFromConfig(cfg *config.Config) ([]SecretRef, error) {
	refs := make([]SecretRef, 0, len(cfg.Secrets))
	for _, s := range cfg.Secrets {
		parsed, err := vault.ParseSecretPath(s.Path)
		if err != nil {
			return nil, fmt.Errorf("invalid secret path %q: %w", s.Path, err)
		}
		refs = append(refs, SecretRef{
			EnvKey: s.Env,
			Path:   parsed.Path,
			Field:  parsed.Field,
		})
	}
	return refs, nil
}
