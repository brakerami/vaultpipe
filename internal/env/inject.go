package env

import (
	"fmt"
	"os"
	"strings"
)

// BuildEnv merges the current process environment with the provided secrets map.
// Secret values take precedence over existing environment variables.
// Returns a slice of "KEY=VALUE" strings ready to pass to exec.Cmd.Env.
func BuildEnv(secrets map[string]string) []string {
	current := os.Environ()
	overridden := make(map[string]bool, len(secrets))

	result := make([]string, 0, len(current)+len(secrets))
	for _, entry := range current {
		key := strings.SplitN(entry, "=", 2)[0]
		if val, ok := secrets[key]; ok {
			result = append(result, fmt.Sprintf("%s=%s", key, val))
			overridden[key] = true
		} else {
			result = append(result, entry)
		}
	}

	// Append secrets that were not already present in the environment
	for k, v := range secrets {
		if !overridden[k] {
			result = append(result, fmt.Sprintf("%s=%s", k, v))
		}
	}
	return result
}

// SanitizeKey converts a secret field name to a valid environment variable key
// by uppercasing and replacing non-alphanumeric characters with underscores.
func SanitizeKey(key string) string {
	upper := strings.ToUpper(key)
	var sb strings.Builder
	for _, r := range upper {
		if (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			sb.WriteRune(r)
		} else {
			sb.WriteRune('_')
		}
	}
	return sb.String()
}
