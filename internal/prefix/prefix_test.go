package prefix_test

import (
	"context"
	"testing"

	"github.com/your-org/vaultpipe/internal/prefix"
)

func TestNew_EmptyPrefix_ReturnsError(t *testing.T) {
	_, err := prefix.New("")
	if err == nil {
		t.Fatal("expected error for empty prefix")
	}
}

func TestNew_ConsecutiveSlashes_ReturnsError(t *testing.T) {
	_, err := prefix.New("secret//data")
	if err == nil {
		t.Fatal("expected error for consecutive slashes")
	}
}

func TestNew_Valid(t *testing.T) {
	p, err := prefix.New("secret/data")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.String() != "secret/data" {
		t.Fatalf("expected 'secret/data', got %q", p.String())
	}
}

func TestApply_JoinsWithSlash(t *testing.T) {
	p, _ := prefix.New("secret/data")
	got := p.Apply("myapp/db")
	if got != "secret/data/myapp/db" {
		t.Fatalf("expected 'secret/data/myapp/db', got %q", got)
	}
}

func TestApply_StripsLeadingSlashFromPath(t *testing.T) {
	p, _ := prefix.New("secret/data")
	got := p.Apply("/myapp/db")
	if got != "secret/data/myapp/db" {
		t.Fatalf("expected 'secret/data/myapp/db', got %q", got)
	}
}

func TestApply_TrailingSlashOnPrefix_Normalised(t *testing.T) {
	p, _ := prefix.New("secret/data/")
	got := p.Apply("key")
	if got != "secret/data/key" {
		t.Fatalf("expected 'secret/data/key', got %q", got)
	}
}

func TestApply_EmptyPath_ReturnsPrefix(t *testing.T) {
	p, _ := prefix.New("secret/data")
	got := p.Apply("")
	if got != "secret/data" {
		t.Fatalf("expected 'secret/data', got %q", got)
	}
}

func TestWrap_PrependsPrefixBeforeDelegate(t *testing.T) {
	p, _ := prefix.New("kv/prod")

	var capturedPath string
	next := func(_ context.Context, path string) (string, error) {
		capturedPath = path
		return "s3cr3t", nil
	}

	wrapped := p.Wrap(next)
	val, err := wrapped(context.Background(), "service/password")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "s3cr3t" {
		t.Fatalf("expected 's3cr3t', got %q", val)
	}
	if capturedPath != "kv/prod/service/password" {
		t.Fatalf("expected 'kv/prod/service/password', got %q", capturedPath)
	}
}
