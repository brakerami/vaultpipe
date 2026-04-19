// Package audit provides structured JSON audit logging for vaultpipe.
//
// It records secret resolution attempts and process lifecycle events,
// writing newline-delimited JSON to a configurable io.Writer so that
// operators can pipe output to log aggregators or files without any
// secrets ever touching disk through this package.
package audit
