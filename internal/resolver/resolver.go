package resolver

import (
	"fmt"

	"github.com/your-org/vaultpipe/internal/vault"
)

// SecretRef holds a mapping from an env var key to a Vault secret path.
type SecretRef struct {
	EnvKey string
	Path   string
	Field  string
}

// Resolver fetches secrets from Vault and returns key=value pairs.
type Resolver struct {
	client *vault.Client
}

// New creates a new Resolver with the given Vault client.
func New(client *vault.Client) *Resolver {
	return &Resolver{client: client}
}

// Resolve takes a slice of SecretRefs and returns a map of envKey -> secretValue.
func (r *Resolver) Resolve(refs []SecretRef) (map[string]string, error) {
	result := make(map[string]string, len(refs))
	for _, ref := range refs {
		value, err := r.client.GetSecret(ref.Path, ref.Field)
		if err != nil {
			return nil, fmt.Errorf("resolving %s from %s#%s: %w", ref.EnvKey, ref.Path, ref.Field, err)
		}
		result[ref.EnvKey] = value
	}
	return result, nil
}
