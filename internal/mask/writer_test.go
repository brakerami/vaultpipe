package mask_test

import (
	"bytes"
	"testing"

	"github.com/yourorg/vaultpipe/internal/mask"
)

func TestWriter_MasksOnWrite(t *testing.T) {
	m := mask.New()
	m.Add("topsecret")

	var buf bytes.Buffer
	w := mask.NewWriter(&buf, m)

	input := []byte("value=topsecret")
	n, err := w.Write(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != len(input) {
		t.Errorf("expected n=%d, got %d", len(input), n)
	}

	got := buf.String()
	if got != "value=[REDACTED]" {
		t.Errorf("got %q", got)
	}
}

func TestWriter_PassthroughWhenNoSecrets(t *testing.T) {
	m := mask.New()
	var buf bytes.Buffer
	w := mask.NewWriter(&buf, m)

	_, _ = w.Write([]byte("plain text"))
	if buf.String() != "plain text" {
		t.Errorf("unexpected output: %q", buf.String())
	}
}
