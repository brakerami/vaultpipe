// Package watch implements TTL-aware secret lease watching for vaultpipe.
//
// A Watcher tracks one or more Vault secret references, each with an
// associated TTL. When a secret approaches expiry (configurable threshold,
// default 25% of TTL remaining), a caller-supplied RenewFunc is invoked so
// the secret can be re-fetched and the environment refreshed without
// restarting the child process.
//
// Usage:
//
//	w := watch.New(myRenewFunc, 0.25)
//	w.Add(ctx, "secret/data/db#password", 30*time.Second)
//	// ... later ...
//	w.Stop()
//
// Helper constructors LoggingRenewFunc and ChannelRenewFunc wrap an
// existing RenewFunc to add observability without changing core logic.
package watch
