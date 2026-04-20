// Package environ provides environment snapshot utilities for vaultpipe.
//
// A Snapshot captures the current process environment at a point in time
// and allows merging injected secrets on top without modifying the live
// os environment. The resulting slice of KEY=VALUE strings can be passed
// directly to exec.Cmd.Env, ensuring secrets never leak back into the
// parent process after the child exits.
package environ
