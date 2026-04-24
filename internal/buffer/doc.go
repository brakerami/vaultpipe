// Package buffer implements a thread-safe, fixed-capacity ring buffer
// designed for capturing recent diagnostic messages, audit events, or
// secret-fetch records within vaultpipe.
//
// The ring overwrites the oldest entry when full, making it suitable for
// bounded memory usage in long-running processes. Use Tail or TailString
// to render the most recent N entries for display or logging.
package buffer
