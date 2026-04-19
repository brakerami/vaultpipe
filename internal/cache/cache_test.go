package cache_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultpipe/internal/cache"
)

func TestCache_SetAndGet(t *testing.T) {
	c := cache.New(5 * time.Second)
	c.Set("mykey", "myvalue")

	v, ok := c.Get("mykey")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if v != "myvalue" {
		t.Fatalf("expected 'myvalue', got %q", v)
	}
}

func TestCache_Miss(t *testing.T) {
	c := cache.New(5 * time.Second)
	_, ok := c.Get("nonexistent")
	if ok {
		t.Fatal("expected cache miss")
	}
}

func TestCache_Expiry(t *testing.T) {
	c := cache.New(10 * time.Millisecond)
	c.Set("k", "v")

	time.Sleep(20 * time.Millisecond)

	_, ok := c.Get("k")
	if ok {
		t.Fatal("expected entry to be expired")
	}
}

func TestCache_ZeroTTL_DisablesCaching(t *testing.T) {
	c := cache.New(0)
	c.Set("k", "v")
	_, ok := c.Get("k")
	if ok {
		t.Fatal("expected cache to be disabled with zero TTL")
	}
}

func TestCache_Invalidate(t *testing.T) {
	c := cache.New(5 * time.Second)
	c.Set("k", "v")
	c.Invalidate("k")
	_, ok := c.Get("k")
	if ok {
		t.Fatal("expected key to be invalidated")
	}
}

func TestCache_Len(t *testing.T) {
	c := cache.New(5 * time.Second)
	if c.Len() != 0 {
		t.Fatal("expected empty cache")
	}
	c.Set("a", "1")
	c.Set("b", "2")
	if c.Len() != 2 {
		t.Fatalf("expected 2 entries, got %d", c.Len())
	}
}
