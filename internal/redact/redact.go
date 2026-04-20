// Package redact provides utilities for scrubbing secret values
// from log output and error messages before they are surfaced to users.
package redact

import "strings"

const placeholder = "***"

// Redactor holds a set of known secret values and replaces them in strings.
type Redactor struct {
	secrets []string
}

// New returns a new Redactor seeded with the provided secret values.
// Empty strings are silently ignored.
func New(secrets []string) *Redactor {
	filtered := make([]string, 0, len(secrets))
	for _, s := range secrets {
		if s != "" {
			filtered = append(filtered, s)
		}
	}
	return &Redactor{secrets: filtered}
}

// Add registers an additional secret value with the Redactor.
func (r *Redactor) Add(secret string) {
	if secret != "" {
		r.secrets = append(r.secrets, secret)
	}
}

// Scrub replaces all known secret values in s with the placeholder.
func (r *Redactor) Scrub(s string) string {
	for _, secret := range r.secrets {
		s = strings.ReplaceAll(s, secret, placeholder)
	}
	return s
}

// ScrubError returns a new error whose message has secrets replaced,
// or nil if err is nil.
func (r *Redactor) ScrubError(err error) error {
	if err == nil {
		return nil
	}
	return &redactedError{msg: r.Scrub(err.Error())}
}

// Len returns the number of secret values currently registered with the Redactor.
func (r *Redactor) Len() int {
	return len(r.secrets)
}

type redactedError struct{ msg string }

func (e *redactedError) Error() string { return e.msg }
