// Package circuit provides a three-state circuit breaker (closed, open,
// half-open) for wrapping calls to external systems such as Vault.
//
// When consecutive failures exceed the configured threshold the breaker
// opens and immediately rejects calls with ErrOpen, preventing repeated
// requests to an unhealthy backend. After the reset timeout the breaker
// moves to half-open, allowing a single probe request through. A
// successful probe closes the circuit; a failed probe re-opens it.
package circuit
