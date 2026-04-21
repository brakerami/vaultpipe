package metrics_test

import (
	"sync"
	"testing"

	"github.com/your-org/vaultpipe/internal/metrics"
)

func TestCounter_StartsAtZero(t *testing.T) {
	r := metrics.New()
	c := r.Counter("test.counter")
	if c.Value() != 0 {
		t.Fatalf("expected 0, got %d", c.Value())
	}
}

func TestCounter_Inc(t *testing.T) {
	r := metrics.New()
	c := r.Counter("test.inc")
	c.Inc()
	c.Inc()
	if got := c.Value(); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestRegistry_SameNameReturnsSameCounter(t *testing.T) {
	r := metrics.New()
	a := r.Counter("shared")
	b := r.Counter("shared")
	a.Inc()
	if b.Value() != 1 {
		t.Fatal("expected same counter instance")
	}
}

func TestRegistry_Snapshot(t *testing.T) {
	r := metrics.New()
	r.Counter("a").Inc()
	r.Counter("a").Inc()
	r.Counter("b").Inc()

	snap := r.Snapshot()
	if snap["a"] != 2 {
		t.Errorf("expected a=2, got %d", snap["a"])
	}
	if snap["b"] != 1 {
		t.Errorf("expected b=1, got %d", snap["b"])
	}
}

func TestRegistry_Reset(t *testing.T) {
	r := metrics.New()
	r.Counter("x").Inc()
	r.Reset()
	if r.Counter("x").Value() != 0 {
		t.Fatal("expected counter reset to 0")
	}
}

func TestCounter_ConcurrentInc(t *testing.T) {
	r := metrics.New()
	c := r.Counter("concurrent")
	const goroutines = 100
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			c.Inc()
		}()
	}
	wg.Wait()
	if got := c.Value(); got != goroutines {
		t.Fatalf("expected %d, got %d", goroutines, got)
	}
}
