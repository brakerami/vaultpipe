package redact

import (
	"bytes"
	"testing"
)

func TestWriter_RedactsOnWrite(t *testing.T) {
	var buf bytes.Buffer
	r := New([]string{"password"})
	w := NewWriter(&buf, r)

	input := []byte("auth with password here")
	n, err := w.Write(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != len(input) {
		t.Errorf("expected n=%d got %d", len(input), n)
	}
	if buf.String() != "auth with *** here" {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestWriter_PassthroughWhenNoSecrets(t *testing.T) {
	var buf bytes.Buffer
	r := New(nil)
	w := NewWriter(&buf, r)

	_, _ = w.Write([]byte("safe text"))
	if buf.String() != "safe text" {
		t.Errorf("unexpected: %q", buf.String())
	}
}
