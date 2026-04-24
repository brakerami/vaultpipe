// Package sieve provides path-based allow/deny filtering for Vault secret
// references. Rules are evaluated in insertion order; the first rule whose
// prefix matches the requested path determines the outcome. If no rule
// matches, the path is permitted by default.
//
// # Basic usage
//
//	s := sieve.New()
//	s.Deny("secret/data/internal")   // block internal paths
//	s.Allow("secret/data/")          // allow everything else under secret/data/
//
//	if err := s.Check(path); err != nil {
//	    return err
//	}
//
// # Middleware
//
// GuardedFetch wraps any FetchFunc so that denied paths are rejected before
// a network call is made:
//
//	guarded := sieve.GuardedFetch(s, client.Fetch)
package sieve
