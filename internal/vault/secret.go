package vault

import (
	"fmt"
	"strings"
)

// SecretPath represents a parsed Vault secret path with optional field selector.
// Format: "secret/data/myapp#field" or "secret/data/myapp" (returns all fields)
type SecretPath struct {
	Path  string
	Field string
}

// ParseSecretPath parses a secret path string into a SecretPath.
func ParseSecretPath(raw string) SecretPath {
	parts := strings.SplitN(raw, "#", 2)
	sp := SecretPath{Path: parts[0]}
	if len(parts) == 2 {
		sp.Field = parts[1]
	}
	return sp
}

// FetchSecrets retrieves secrets from Vault at the given path and returns
// a map of key=value pairs suitable for injection into a process environment.
func (c *Client) FetchSecrets(sp SecretPath) (map[string]string, error) {
	secret, err := c.Logical().Read(sp.Path)
	if err != nil {
		return nil, fmt.Errorf("reading secret at %q: %w", sp.Path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("no secret found at %q", sp.Path)
	}

	// KV v2 wraps data under secret.Data["data"]
	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		// KV v1 — data is at top level
		data = make(map[string]interface{}, len(secret.Data))
		for k, v := range secret.Data {
			data[k] = v
		}
	}

	if sp.Field != "" {
		v, exists := data[sp.Field]
		if !exists {
			return nil, fmt.Errorf("field %q not found in secret at %q", sp.Field, sp.Path)
		}
		return map[string]string{sp.Field: fmt.Sprintf("%v", v)}, nil
	}

	result := make(map[string]string, len(data))
	for k, v := range data {
		result[k] = fmt.Sprintf("%v", v)
	}
	return result, nil
}
