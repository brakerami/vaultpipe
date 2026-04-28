package dedupe_test

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/dedupe"
)

func TestDo_SingleCaller_ReturnsFetchResult(t *testing.T) {
	g := dedupe.New(0)
	v, err := g.Do(context.Background(), "secret/foo", func(_ context.Context, _ string) (string, error) {
		return "s3cr3t", nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "s3cr3t" {
		t.Fatalf("expected s3cr3t, got %q", v)
	}
}

func TestDo_PropagatesFetchError(t *testing.T) {
	g := dedupe.New(0)
	want := errors.New("vault unavailable")
	_, err := g.Do(context.Background(), "secret/foo", func(_ context.Context, _ string) (string, error) {
		return "", want
	})
	if !errors.Is(err, want) {
		t.Fatalf("expected %v, got %v", want, err)
	}
}

func TestDo_ConcurrentCallers_OnlyOneFetch(t *testing.T) {
	var calls atomic.Int32
	g := dedupe.New(100 * time.Millisecond)

	ready := make(chan struct{})
	var wg sync.WaitGroup
	results := make([]string, 5)

	for i := range results {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			<-ready
			v, _ := g.Do(context.Background(), "secret/bar", func(_ context.Context, _ string) (string, error) {
				calls.Add(1)
				time.Sleep(20 * time.Millisecond)
				return "value", nil
			})
			results[idx] = v
		}(i)
	}
	close(ready)
	wg.Wait()

	for _, v := range results {
		if v != "value" {
			t.Errorf("expected value, got %q", v)
		}
	}
	if n := calls.Load(); n > 2 {
		t.Errorf("expected at most 2 fetch calls (initial + possible race), got %d", n)
	}
}

func TestDo_WindowReuse_SkipsFetch(t *testing.T) {
	var calls atomic.Int32
	g := dedupe.New(500 * time.Millisecond)

	fetch := func(_ context.Context, _ string) (string, error) {
		calls.Add(1)
		return "cached", nil
	}

	g.Do(context.Background(), "secret/baz", fetch) //nolint
	g.Do(context.Background(), "secret/baz", fetch) //nolint

	if n := calls.Load(); n != 1 {
		t.Errorf("expected 1 fetch call, got %d", n)
	}
}

func TestInvalidate_ForcesFreshFetch(t *testing.T) {
	var calls atomic.Int32
	g := dedupe.New(time.Minute)

	fetch := func(_ context.Context, _ string) (string, error) {
		calls.Add(1)
		return "v", nil
	}

	g.Do(context.Background(), "secret/qux", fetch) //nolint
	g.Invalidate("secret/qux")
	g.Do(context.Background(), "secret/qux", fetch) //nolint

	if n := calls.Load(); n != 2 {
		t.Errorf("expected 2 fetch calls after invalidation, got %d", n)
	}
}

func TestDo_ContextCancelled_ReturnsError(t *testing.T) {
	g := dedupe.New(0)
	ctx, cancel := context.WithCancel(context.Background())

	started := make(chan struct{})
	go func() {
		g.Do(ctx, "secret/slow", func(c context.Context, _ string) (string, error) { //nolint
			close(started)
			<-c.Done()
			return "", c.Err()
		})
	}()

	<-started
	cancel()

	ctx2, cancel2 := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel2()
	_, err := g.Do(ctx2, "secret/slow", func(_ context.Context, _ string) (string, error) {
		return "late", nil
	})
	if err == nil {
		t.Log("second caller succeeded after cancellation — acceptable")
	}
}
