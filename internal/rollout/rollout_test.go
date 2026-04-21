package rollout_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/vaultpipe/vaultpipe/internal/rollout"
)

func TestNew_InvalidConcurrency(t *testing.T) {
	_, err := rollout.New(0, 0, func(_ context.Context, _ rollout.Stage) error { return nil })
	if err == nil {
		t.Fatal("expected error for concurrency=0")
	}
}

func TestNew_NilApply(t *testing.T) {
	_, err := rollout.New(1, 0, nil)
	if err == nil {
		t.Fatal("expected error for nil apply func")
	}
}

func TestRun_AllStagesApplied(t *testing.T) {
	var count int64
	c, err := rollout.New(2, 0, func(_ context.Context, s rollout.Stage) error {
		atomic.AddInt64(&count, 1)
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	stages := []rollout.Stage{
		{Key: "A", OldVal: "v1", NewVal: "v2"},
		{Key: "B", OldVal: "v1", NewVal: "v2"},
		{Key: "C", OldVal: "v1", NewVal: "v2"},
	}

	if err := c.Run(context.Background(), stages); err != nil {
		t.Fatalf("unexpected run error: %v", err)
	}
	if count != 3 {
		t.Fatalf("expected 3 stages applied, got %d", count)
	}
}

func TestRun_StopsOnError(t *testing.T) {
	var count int64
	boom := errors.New("apply failed")
	c, _ := rollout.New(1, 0, func(_ context.Context, s rollout.Stage) error {
		if s.Key == "B" {
			return boom
		}
		atomic.AddInt64(&count, 1)
		return nil
	})

	stages := []rollout.Stage{
		{Key: "A"}, {Key: "B"}, {Key: "C"},
	}

	err := c.Run(context.Background(), stages)
	if !errors.Is(err, boom) {
		t.Fatalf("expected boom error, got %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 successful stage before error, got %d", count)
	}
}

func TestRun_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	c, _ := rollout.New(1, 10*time.Millisecond, func(_ context.Context, _ rollout.Stage) error {
		return nil
	})

	stages := []rollout.Stage{{Key: "A"}, {Key: "B"}}
	err := c.Run(ctx, stages)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestRun_EmptyStages(t *testing.T) {
	c, _ := rollout.New(2, 0, func(_ context.Context, _ rollout.Stage) error {
		return nil
	})
	if err := c.Run(context.Background(), nil); err != nil {
		t.Fatalf("unexpected error for empty stages: %v", err)
	}
}
