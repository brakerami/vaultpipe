// Package mask provides secret-aware redaction utilities for vaultpipe.
//
// A Masker holds a registry of known secret values. Any string or byte
// stream can be filtered through the Masker so that secret values are
// replaced with [REDACTED] before reaching logs, stdout, or stderr.
//
// Use NewWriter to wrap an io.Writer (e.g. os.Stderr) so that process
// output is automatically scrubbed before display.
package mask
