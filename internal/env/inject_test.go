package env

import (
	"os"
	"strings"
	"testing"
)

func TestBuildEnv_InjectsNewKeys(t *testing.T) {
	secrets := map[string]string{"MY_SECRET": "supersecret"}
	env := BuildEnv(secrets)

	for _, entry := range env {
		if entry == "MY_SECRET=supersecret" {
			return
		}
	}
	t.Error("expected MY_SECRET=supersecret in env")
}

func TestBuildEnv_OverridesExisting(t *testing.T) {
	os.Setenv("OVERRIDE_ME", "original")
	defer os.Unsetenv("OVERRIDE_ME")

	secrets := map[string]string{"OVERRIDE_ME": "replaced"}
	env := BuildEnv(secrets)

	count := 0
	for _, entry := range env {
		if strings.HasPrefix(entry, "OVERRIDE_ME=") {
			count++
			if entry != "OVERRIDE_ME=replaced" {
				t.Errorf("expected OVERRIDE_ME=replaced, got %q", entry)
			}
		}
	}
	if count != 1 {
		t.Errorf("expected exactly 1 OVERRIDE_ME entry, got %d", count)
	}
}

func TestSanitizeKey(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"db-password", "DB_PASSWORD"},
		{"api.key", "API_KEY"},
		{"MY_VAR", "MY_VAR"},
		{"hello world", "HELLO_WORLD"},
	}
	for _, tc := range cases {
		got := SanitizeKey(tc.input)
		if got != tc.expected {
			t.Errorf("SanitizeKey(%q) = %q, want %q", tc.input, got, tc.expected)
		}
	}
}
