// Package rollout provides a staged secret rollout controller for vaultpipe.
//
// When a secret changes value, multiple consumers may need to be updated in a
// controlled fashion rather than all at once. The Controller batches consumers
// into groups of configurable concurrency, applies each batch in parallel, and
// inserts an optional delay between batches to allow services time to stabilise
// before the next wave is updated.
//
// Example usage:
//
//	ctrl, err := rollout.New(3, 5*time.Second, func(ctx context.Context, s rollout.Stage) error {
//		// restart or signal the consumer identified by s.Key
//		return notifyConsumer(ctx, s.Key, s.NewVal)
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	if err := ctrl.Run(ctx, stages); err != nil {
//		log.Printf("rollout failed: %v", err)
//	}
package rollout
