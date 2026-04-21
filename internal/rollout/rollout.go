// Package rollout provides a staged secret rollout controller that
// coordinates secret rotation across multiple consumers with configurable
// concurrency and per-stage hooks.
package rollout

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Stage represents a single rollout step applied to one consumer.
type Stage struct {
	Key    string
	OldVal string
	NewVal string
}

// ApplyFunc is called for each stage during a rollout.
type ApplyFunc func(ctx context.Context, s Stage) error

// Controller manages staged secret rollouts.
type Controller struct {
	mu          sync.Mutex
	concurrency int
	delay       time.Duration
	apply       ApplyFunc
}

// New creates a Controller. concurrency controls how many stages run in
// parallel; delay is the pause between batches.
func New(concurrency int, delay time.Duration, apply ApplyFunc) (*Controller, error) {
	if concurrency < 1 {
		return nil, fmt.Errorf("rollout: concurrency must be >= 1, got %d", concurrency)
	}
	if apply == nil {
		return nil, fmt.Errorf("rollout: apply func must not be nil")
	}
	return &Controller{
		concurrency: concurrency,
		delay:       delay,
		apply:       apply,
	}, nil
}

// Run executes stages in batches of c.concurrency, pausing c.delay between
// batches. It stops early and returns the first error encountered.
func (c *Controller) Run(ctx context.Context, stages []Stage) error {
	for i := 0; i < len(stages); i += c.concurrency {
		end := i + c.concurrency
		if end > len(stages) {
			end = len(stages)
		}
		batch := stages[i:end]

		if err := c.runBatch(ctx, batch); err != nil {
			return err
		}

		if end < len(stages) && c.delay > 0 {
			select {
			case <-time.After(c.delay):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
	return nil
}

func (c *Controller) runBatch(ctx context.Context, batch []Stage) error {
	type result struct{ err error }
	results := make(chan result, len(batch))

	for _, s := range batch {
		s := s
		go func() {
			results <- result{err: c.apply(ctx, s)}
		}()
	}

	for range batch {
		r := <-results
		if r.err != nil {
			return r.err
		}
	}
	return nil
}
