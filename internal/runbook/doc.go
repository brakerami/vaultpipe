// Package runbook implements pre-flight validation for vaultpipe.
//
// Before spawning the target process, vaultpipe can optionally run a set of
// checks to verify that:
//
//  1. The configured Vault server is reachable and the token is valid.
//  2. Every secret path referenced in the configuration can be read with the
//     current token's policy.
//
// Checks are executed concurrently where possible and results are collected
// into a []CheckResult slice so callers can decide whether to abort or
// proceed with degraded configuration.
//
// Usage:
//
//	runner := runbook.New(vaultClient, cfg)
//	results, err := runner.Run(ctx)
//	if summary, ok := runbook.Summary(results); !ok {
//	    log.Fatal(summary)
//	}
package runbook
