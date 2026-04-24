// Package stagger provides a mechanism to spread concurrent Vault secret
// renewals or fetches across a configurable time window. When many leases
// expire at the same time, stagger prevents a thundering-herd of simultaneous
// requests hitting Vault by introducing a per-caller random delay before
// the operation is executed.
//
// Example usage:
//
//	s := stagger.New(30 * time.Second)
//	if err := s.Do(ctx, renewFn); err != nil {
//		log.Printf("renewal skipped: %v", err)
//	}
package stagger
