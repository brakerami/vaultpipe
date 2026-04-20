// Package health implements a lightweight Vault health probe.
//
// It queries the /v1/sys/health endpoint and interprets HTTP status
// codes according to Vault's documented conventions:
//
//   - 200: initialised, unsealed, and active
//   - 429: unsealed and in standby mode
//   - 503: sealed
//
// The resulting [Status] value is used by the runbook diagnostics
// subsystem to surface actionable information before secret resolution
// begins, and by the startup validation path to fail fast when Vault
// is unavailable.
package health
