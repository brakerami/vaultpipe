package drain_test

import (
	"context"
	"testing"
	"time"

	"github.com/your-org/vaultpipe/internal/drain"
)

func TestNew_DefaultTimeout(t *testing.T) {
	d := drain.New(0)
	if d == nil {
		t.Fatal("expected non-nil Drainer")
	}
}

func TestAcquire_BeforeDrain(t *testing.T) {
	d := drain.New(time.Second)
	if !d.Acquire() {
		t.Fatal("expected Acquire to return true before Drain")
	}
	d.Release()
}

func TestAcquire_AfterDrain_ReturnsFalse(t *testing.T) {
	d := drain.New(time.Second)

	// Drain with nothing acquired — should return immediately.
	if err := d.Drain(context.Background()); err != nil {
		t.Fatalf("unexpected drain error: %v", err)
	}

	if d.Acquire() {
		t.Fatal("expected Acquire to return false after Drain")
	}
}

func TestDrain_WaitsForRelease(t *testing.T) {
	d := drain.New(time.Second)

	if !d.Acquire() {
		t.Fatal("expected Acquire to succeed")
	}

	released := make(chan struct{})
	go func() {
		time.Sleep(20 * time.Millisecond)
		d.Release()
		close(released)
	}()

	start := time.Now()
	if err := d.Drain(context.Background()); err != nil {
		t.Fatalf("unexpected drain error: %v", err)
	}

	<-released
	if time.Since(start) < 20*time.Millisecond {
		t.Error("Drain returned before Release was called")
	}
}

func TestDrain_TimesOut(t *testing.T) {
	d := drain.New(30 * time.Millisecond)

	if !d.Acquire() {
		t.Fatal("expected Acquire to succeed")
	}
	defer d.Release() // never released in time — intentional

	err := d.Drain(context.Background())
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}

func TestDrain_ContextCancelled(t *testing.T) {
	d := drain.New(5 * time.Second)

	if !d.Acquire() {
		t.Fatal("expected Acquire to succeed")
	}
	defer d.Release()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	err := d.Drain(ctx)
	if err == nil {
		t.Fatal("expected context error, got nil")
	}
}
