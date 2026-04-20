// Package template provides simple secret reference expansion
// for environment variable values containing {{vault:path#field}} placeholders.
package template

import (
	"fmt"
	"regexp"
	"strings"
)

// refPattern matches {{vault:secret/path#field}} placeholders.
var refPattern = regexp.MustCompile(`\{\{vault:([^}]+)\}\}`)

// Resolver resolves a secret path+field to a value.
type Resolver interface {
	Resolve(path, field string) (string, error)
}

// Expand replaces all {{vault:path#field}} placeholders in s with resolved
// secret values. Returns an error if any placeholder cannot be resolved.
func Expand(s string, r Resolver) (string, error) {
	var expandErr error
	result := refPattern.ReplaceAllStringFunc(s, func(match string) string {
		if expandErr != nil {
			return match
		}
		// Extract inner: path#field
		inner := match[8 : len(match)-2] // strip "{{vault:" and "}}"
		path, field, err := splitRef(inner)
		if err != nil {
			expandErr = fmt.Errorf("template: invalid ref %q: %w", match, err)
			return match
		}
		val, err := r.Resolve(path, field)
		if err != nil {
			expandErr = fmt.Errorf("template: resolve %q: %w", match, err)
			return match
		}
		return val
	})
	if expandErr != nil {
		return "", expandErr
	}
	return result, nil
}

// HasRefs reports whether s contains any vault placeholder references.
func HasRefs(s string) bool {
	return refPattern.MatchString(s)
}

func splitRef(ref string) (path, field string, err error) {
	idx := strings.LastIndex(ref, "#")
	if idx < 0 {
		return "", "", fmt.Errorf("missing '#field' in %q", ref)
	}
	path = ref[:idx]
	field = ref[idx+1:]
	if path == "" || field == "" {
		return "", "", fmt.Errorf("path and field must be non-empty in %q", ref)
	}
	return path, field, nil
}
