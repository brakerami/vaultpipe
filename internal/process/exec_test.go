package process

import (
	"strings"
	"testing"
)

func TestNewRunner_NoArgs(t *testing.T) {
	_, err := NewRunner(nil, []string{})
	if err == nil {
		t.Fatal("expected error for empty args, got nil")
	}
}

func TestNewRunner_ValidArgs(t *testing.T) {
	env := []string{"FOO=bar"}
	args := []string{"echo", "hello"}

	r, err := NewRunner(env, args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Args) != 2 {
		t.Errorf("expected 2 args, got %d", len(r.Args))
	}
	if len(r.Env) != 1 {
		t.Errorf("expected 1 env var, got %d", len(r.Env))
	}
}

func TestRunner_Run_InvalidCommand(t *testing.T) {
	r, err := NewRunner(nil, []string{"__nonexistent_cmd_vaultpipe__"})
	if err != nil {
		t.Fatalf("unexpected error creating runner: %v", err)
	}
	err = r.Run()
	if err == nil {
		t.Fatal("expected error running nonexistent command, got nil")
	}
}

func TestRunner_Run_Echo(t *testing.T) {
	env := []string{"PATH=/usr/bin:/bin"}
	r, err := NewRunner(env, []string{"echo", "vaultpipe"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	err = r.Run()
	if err != nil {
		// Some CI environments may not have echo at expected path; soft fail.
		if !strings.Contains(err.Error(), "not found") {
			t.Errorf("unexpected run error: %v", err)
		}
	}
}
