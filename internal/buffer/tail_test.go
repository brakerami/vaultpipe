package buffer_test

import (
	"strings"
	"testing"

	"github.com/yourusername/vaultpipe/internal/buffer"
)

func TestTail_AllEntries(t *testing.T) {
	r := buffer.New(10)
	r.Add("alpha")
	r.Add("beta")
	r.Add("gamma")

	out := buffer.TailString(r, 0)
	for _, want := range []string{"alpha", "beta", "gamma"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output: %q", want, out)
		}
	}
}

func TestTail_LimitsToN(t *testing.T) {
	r := buffer.New(10)
	r.Add("one")
	r.Add("two")
	r.Add("three")

	out := buffer.TailString(r, 2)
	if strings.Contains(out, "one") {
		t.Error("expected 'one' to be excluded from tail(2)")
	}
	if !strings.Contains(out, "two") || !strings.Contains(out, "three") {
		t.Errorf("expected 'two' and 'three' in tail(2): %q", out)
	}
}

func TestTail_EmptyRing(t *testing.T) {
	r := buffer.New(5)
	out := buffer.TailString(r, 0)
	if out != "" {
		t.Errorf("expected empty string for empty ring, got %q", out)
	}
}

func TestTail_NGreaterThanLen(t *testing.T) {
	r := buffer.New(5)
	r.Add("only")
	out := buffer.TailString(r, 100)
	if !strings.Contains(out, "only") {
		t.Errorf("expected 'only' in output: %q", out)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 1 {
		t.Errorf("expected 1 line, got %d", len(lines))
	}
}

func TestTail_WriterError(t *testing.T) {
	r := buffer.New(4)
	r.Add("msg")
	ew := &errWriter{}
	n, err := buffer.Tail(ew, r, 0)
	if err == nil {
		t.Error("expected error from failing writer")
	}
	if n != 0 {
		t.Errorf("expected 0 written, got %d", n)
	}
}

type errWriter struct{}

func (e *errWriter) Write(_ []byte) (int, error) {
	return 0, fmt.Errorf("write error")
}

import "fmt"
