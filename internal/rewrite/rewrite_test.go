package rewrite_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/rewrite"
)

func TestNew_UnknownKind(t *testing.T) {
	_, err := rewrite.New([]rewrite.Rule{{Kind: "rot13"}})
	if err == nil {
		t.Fatal("expected error for unknown kind")
	}
}

func TestNew_PrefixRequiresArg(t *testing.T) {
	_, err := rewrite.New([]rewrite.Rule{{Kind: "prefix", Arg: ""}})
	if err == nil {
		t.Fatal("expected error when prefix arg is empty")
	}
}

func TestNew_SuffixRequiresArg(t *testing.T) {
	_, err := rewrite.New([]rewrite.Rule{{Kind: "suffix", Arg: ""}})
	if err == nil {
		t.Fatal("expected error when suffix arg is empty")
	}
}

func TestApply_Upper(t *testing.T) {
	rw, _ := rewrite.New([]rewrite.Rule{{Kind: "upper"}})
	if got := rw.Apply("hello"); got != "HELLO" {
		t.Fatalf("want HELLO, got %q", got)
	}
}

func TestApply_Lower(t *testing.T) {
	rw, _ := rewrite.New([]rewrite.Rule{{Kind: "lower"}})
	if got := rw.Apply("WORLD"); got != "world" {
		t.Fatalf("want world, got %q", got)
	}
}

func TestApply_Trim(t *testing.T) {
	rw, _ := rewrite.New([]rewrite.Rule{{Kind: "trim"}})
	if got := rw.Apply("  secret  "); got != "secret" {
		t.Fatalf("want 'secret', got %q", got)
	}
}

func TestApply_Prefix(t *testing.T) {
	rw, _ := rewrite.New([]rewrite.Rule{{Kind: "prefix", Arg: "prod_"}})
	if got := rw.Apply("token"); got != "prod_token" {
		t.Fatalf("want prod_token, got %q", got)
	}
}

func TestApply_Suffix(t *testing.T) {
	rw, _ := rewrite.New([]rewrite.Rule{{Kind: "suffix", Arg: "_v2"}})
	if got := rw.Apply("key"); got != "key_v2" {
		t.Fatalf("want key_v2, got %q", got)
	}
}

func TestApply_ChainedRules(t *testing.T) {
	rw, _ := rewrite.New([]rewrite.Rule{
		{Kind: "trim"},
		{Kind: "upper"},
		{Kind: "prefix", Arg: "[["},
		{Kind: "suffix", Arg: "]]"},
	})
	if got := rw.Apply("  hello  "); got != "[[HELLO]]" {
		t.Fatalf("want [[HELLO]], got %q", got)
	}
}

func TestApply_NoRules_Passthrough(t *testing.T) {
	rw, _ := rewrite.New(nil)
	if got := rw.Apply("unchanged"); got != "unchanged" {
		t.Fatalf("want unchanged, got %q", got)
	}
}

func TestLen(t *testing.T) {
	rw, _ := rewrite.New([]rewrite.Rule{{Kind: "upper"}, {Kind: "trim"}})
	if rw.Len() != 2 {
		t.Fatalf("want 2, got %d", rw.Len())
	}
}
