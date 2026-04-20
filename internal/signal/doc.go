// Package signal provides transparent OS signal forwarding from the vaultpipe
// parent process to the spawned child process.
//
// This ensures that signals such as SIGINT, SIGTERM, and SIGHUP are correctly
// propagated so the child process can perform graceful shutdown without
// vaultpipe intercepting or swallowing them.
package signal
