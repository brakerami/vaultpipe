// Package redact provides a Redactor type that scrubs registered secret
// values from strings, errors, and io.Writer streams before they are
// written to logs or returned to callers.
//
// Use New to create a Redactor, Add to register secrets at runtime, and
// NewWriter to wrap any io.Writer with automatic scrubbing.
package redact
