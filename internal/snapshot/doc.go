// Package snapshot provides point-in-time capture of resolved secret
// environments and utilities for computing diffs between two snapshots.
//
// Typical usage:
//
//	prev := snapshot.Take(oldSecrets)
//	next := snapshot.Take(newSecrets)
//	changes := snapshot.Diff(prev, next)
//	for _, c := range changes {
//		fmt.Printf("%s %s\n", c.Kind, c.Key)
//	}
package snapshot
