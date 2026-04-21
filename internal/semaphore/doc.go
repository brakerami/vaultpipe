// Package semaphore implements a counting semaphore for bounding concurrent
// access to shared resources within vaultpipe.
//
// It is used to cap the number of in-flight Vault secret fetches so that
// bursts of process starts do not overwhelm the Vault server.
//
// Example:
//
//	sem, err := semaphore.New(5)
//	if err != nil { ... }
//
//	if err := sem.Acquire(ctx); err != nil {
//		return err
//	}
//	defer sem.Release()
//	// ... perform Vault fetch ...
package semaphore
