package watch_test

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/your-org/vaultpipe/internal/watch"
)

func TestLoggingRenewFunc_WritesOnSuccess(t *testing.T) {
	var buf bytes.Buffer
	inner := func(_ context.Context, _ string) error { return nil }
	fn := watch.LoggingRenewFunc(inner, &buf)

	_ = fn(context.Background(), "secret/data/foo")

	if !strings.Contains(buf.String(), "secret/data/foo") {
		t.Errorf("expected ref in log output, got: %s", buf.String())
	}
	if !strings.Contains(buf.String(), "renewed") {
		t.Errorf("expected 'renewed' in log output, got: %s", buf.String())
	}
}

func TestLoggingRenewFunc_WritesOnError(t *testing.T) {
	var buf bytes.Buffer
	inner := func(_ context.Context, _ string) error { return errors.New("vault down") }
	fn := watch.LoggingRenewFunc(inner, &buf)

	err := fn(context.Background(), "secret/data/bar")
	if err == nil {
		t.Fatal("expected error to propagate")
	}
	if !strings.Contains(buf.String(), "vault down") {
		t.Errorf("expected error message in log output, got: %s", buf.String())
	}
}

func TestChannelRenewFunc_SendsEvent(t *testing.T) {
	ch := make(chan watch.Event, 1)
	inner := func(_ context.Context, _ string) error { return nil }
	fn := watch.ChannelRenewFunc(inner, ch)

	_ = fn(context.Background(), "secret/data/baz")

	select {
	case ev := <-ch:
		if ev.Ref != "secret/data/baz" {
			t.Errorf("unexpected ref: %s", ev.Ref)
		}
		if ev.Err != nil {
			t.Errorf("unexpected error: %v", ev.Err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timed out waiting for event")
	}
}

func TestChannelRenewFunc_NonBlocking(t *testing.T) {
	// Unbuffered — send must not block.
	ch := make(chan watch.Event) // intentionally unbuffered
	inner := func(_ context.Context, _ string) error { return nil }
	fn := watch.ChannelRenewFunc(inner, ch)

	done := make(chan struct{})
	go func() {
		_ = fn(context.Background(), "secret/nb")
		close(done)
	}()

	select {
	case <-done:
		// good — did not block
	case <-time.After(200 * time.Millisecond):
		t.Fatal("ChannelRenewFunc blocked on full channel")
	}
}
