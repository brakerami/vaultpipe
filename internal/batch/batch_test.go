package batch_test

import (
	"context"
	"errors"
	"sort"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/batch"
)

func TestResolve_AllSucceed(t *testing.T) {
	fetch := func(_ context.Context, path string) (string, error) {
		return "val:" + path, nil
	}
	r := batch.New(fetch, 4)
	refs := map[string]string{"A": "secret/a", "B": "secret/b", "C": "secret/c"}

	results := r.Resolve(context.Background(), refs)
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	for _, res := range results {
		if res.Err != nil {
			t.Errorf("unexpected error for key %s: %v", res.Key, res.Err)
		}
		if res.Value != "val:"+refs[res.Key] {
			t.Errorf("key %s: want %q got %q", res.Key, "val:"+refs[res.Key], res.Value)
		}
	}
}

func TestResolve_PartialError(t *testing.T) {
	errBoom := errors.New("boom")
	fetch := func(_ context.Context, path string) (string, error) {
		if path == "secret/bad" {
			return "", errBoom
		}
		return "ok", nil
	}
	r := batch.New(fetch, 2)
	refs := map[string]string{"GOOD": "secret/good", "BAD": "secret/bad"}

	results := r.Resolve(context.Background(), refs)
	var errs int
	for _, res := range results {
		if res.Err != nil {
			errs++
		}
	}
	if errs != 1 {
		t.Fatalf("expected 1 error, got %d", errs)
	}
}

func TestResolve_ConcurrencyRespected(t *testing.T) {
	var active int64
	var peak int64

	fetch := func(_ context.Context, _ string) (string, error) {
		cur := atomic.AddInt64(&active, 1)
		for {
			p := atomic.LoadInt64(&peak)
			if cur <= p || atomic.CompareAndSwapInt64(&peak, p, cur) {
				break
			}
		}
		time.Sleep(10 * time.Millisecond)
		atomic.AddInt64(&active, -1)
		return "v", nil
	}

	refs := make(map[string]string)
	for i := 0; i < 10; i++ {
		refs[string(rune('A'+i))] = "secret/x"
	}

	r := batch.New(fetch, 3)
	r.Resolve(context.Background(), refs)

	if peak > 3 {
		t.Errorf("peak concurrency %d exceeded cap of 3", peak)
	}
}

func TestResolve_EmptyRefs(t *testing.T) {
	r := batch.New(func(_ context.Context, _ string) (string, error) { return "", nil }, 2)
	results := r.Resolve(context.Background(), map[string]string{})
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestResolve_ZeroConcurrency_DefaultsToOne(t *testing.T) {
	fetch := func(_ context.Context, _ string) (string, error) { return "v", nil }
	r := batch.New(fetch, 0)
	refs := map[string]string{"K": "secret/k"}
	results := r.Resolve(context.Background(), refs)
	if len(results) != 1 {
		t.Fatalf("expected 1 result")
	}
	_ = sort.Search // suppress import warning
}
