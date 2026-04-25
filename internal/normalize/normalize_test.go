package normalize_test

import (
	"testing"

	"github.com/your-org/vaultpipe/internal/normalize"
)

func TestKey_AlreadyUpper(t *testing.T) {
	if got := normalize.Key("HELLO"); got != "HELLO" {
		t.Fatalf("expected HELLO, got %q", got)
	}
}

func TestKey_LowerConverted(t *testing.T) {
	if got := normalize.Key("hello"); got != "HELLO" {
		t.Fatalf("expected HELLO, got %q", got)
	}
}

func TestKey_HyphenReplaced(t *testing.T) {
	if got := normalize.Key("my-secret-key"); got != "MY_SECRET_KEY" {
		t.Fatalf("expected MY_SECRET_KEY, got %q", got)
	}
}

func TestKey_DotReplaced(t *testing.T) {
	if got := normalize.Key("db.password"); got != "DB_PASSWORD" {
		t.Fatalf("expected DB_PASSWORD, got %q", got)
	}
}

func TestKey_ConsecutiveSeparatorsCollapsed(t *testing.T) {
	if got := normalize.Key("foo--bar"); got != "FOO_BAR" {
		t.Fatalf("expected FOO_BAR, got %q", got)
	}
}

func TestKey_LeadingTrailingSeparatorsTrimmed(t *testing.T) {
	if got := normalize.Key("-leading-trailing-"); got != "LEADING_TRAILING" {
		t.Fatalf("expected LEADING_TRAILING, got %q", got)
	}
}

func TestKey_WithPrefix(t *testing.T) {
	got := normalize.Key("token", normalize.WithPrefix("APP"))
	if got != "APP_TOKEN" {
		t.Fatalf("expected APP_TOKEN, got %q", got)
	}
}

func TestKey_WithCustomSeparator(t *testing.T) {
	got := normalize.Key("my-key", normalize.WithSeparator('.'))
	if got != "MY.KEY" {
		t.Fatalf("expected MY.KEY, got %q", got)
	}
}

func TestKey_EmptyString(t *testing.T) {
	if got := normalize.Key(""); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestMap_NormalizesAllKeys(t *testing.T) {
	input := map[string]string{
		"db-host":  "localhost",
		"db.port":  "5432",
		"api_key":  "secret",
	}
	out := normalize.Map(input)

	cases := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"API_KEY": "secret",
	}
	for k, want := range cases {
		if got, ok := out[k]; !ok {
			t.Errorf("key %q missing from output", k)
		} else if got != want {
			t.Errorf("key %q: want %q, got %q", k, want, got)
		}
	}
}

func TestMap_WithPrefix(t *testing.T) {
	input := map[string]string{"token": "abc"}
	out := normalize.Map(input, normalize.WithPrefix("VAULT"))
	if v, ok := out["VAULT_TOKEN"]; !ok || v != "abc" {
		t.Fatalf("expected VAULT_TOKEN=abc, got map %v", out)
	}
}
