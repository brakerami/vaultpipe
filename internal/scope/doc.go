// Package scope provides key-prefix scoping for secret environment
// variables. A Scope wraps a namespace prefix and applies it
// consistently when secrets are injected into the process environment,
// preventing collisions between secrets originating from different
// Vault paths or services.
//
// Example usage:
//
//	s, err := scope.New("PAYMENTS", "_")
//	if err != nil { ... }
//	scoped := s.Map(resolvedSecrets) // {"PAYMENTS_DB_PASS": "...", ...}
package scope
