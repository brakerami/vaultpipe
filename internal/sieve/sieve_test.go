package sieve_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/sieve"
)

func TestPermitted_NoRules_AllowsAll(t *testing.T) {
	s := sieve.New()
	if !s.Permitted("secret/data/foo") {
		t.Fatal("expected path to be permitted when no rules are set")
	}
}

func TestPermitted_AllowRule_Matches(t *testing.T) {
	s := sieve.New()
	s.Allow("secret/data/")
	if !s.Permitted("secret/data/myapp") {
		t.Fatal("expected allow rule to permit matching path")
	}
}

func TestPermitted_DenyRule_Blocks(t *testing.T) {
	s := sieve.New()
	s.Deny("secret/data/internal")
	if s.Permitted("secret/data/internal/creds") {
		t.Fatal("expected deny rule to block matching path")
	}
}

func TestPermitted_FirstRuleWins(t *testing.T) {
	s := sieve.New()
	s.Allow("secret/")
	s.Deny("secret/data/internal")
	// allow comes first, so internal path should still be permitted
	if !s.Permitted("secret/data/internal/creds") {
		t.Fatal("expected first matching rule (allow) to win")
	}
}

func TestPermitted_DenyBeforeAllow(t *testing.T) {
	s := sieve.New()
	s.Deny("secret/data/internal")
	s.Allow("secret/")
	if s.Permitted("secret/data/internal/creds") {
		t.Fatal("expected deny rule (first) to block path")
	}
}

func TestCheck_ReturnsErrorWhenDenied(t *testing.T) {
	s := sieve.New()
	s.Deny("forbidden/")
	if err := s.Check("forbidden/path"); err == nil {
		t.Fatal("expected error for denied path")
	}
}

func TestCheck_ReturnsNilWhenAllowed(t *testing.T) {
	s := sieve.New()
	s.Allow("allowed/")
	if err := s.Check("allowed/path"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRules_ReturnsCopy(t *testing.T) {
	s := sieve.New()
	s.Allow("a/")
	s.Deny("b/")
	rules := s.Rules()
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	// mutating returned slice must not affect sieve
	rules[0].Prefix = "mutated/"
	if s.Rules()[0].Prefix != "a/" {
		t.Fatal("Rules() must return an isolated copy")
	}
}

func TestReset_ClearsRules(t *testing.T) {
	s := sieve.New()
	s.Deny("secret/")
	s.Reset()
	if len(s.Rules()) != 0 {
		t.Fatal("expected rules to be empty after Reset")
	}
	if !s.Permitted("secret/data/foo") {
		t.Fatal("expected all paths permitted after Reset")
	}
}
