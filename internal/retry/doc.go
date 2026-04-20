// Package retry implements exponential backoff retry logic for use when
// communicating with HashiCorp Vault. It supports configurable attempt limits,
// base delay, maximum delay, and a multiplier for backoff growth.
//
// Errors can be wrapped with [Permanent] to prevent further retries when the
// failure is not transient (e.g. permission denied, secret not found).
//
// Example:
//
//	cfg := retry.DefaultConfig()
//	err := retry.Do(ctx, cfg, func() error {
//		return vaultClient.Read(path)
//	})
package retry
