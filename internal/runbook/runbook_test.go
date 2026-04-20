package runbook_test

import (
	"context"
	"errors"
	"testing"

	"github.com/yourusername/vaultpipe/internal/config"
	"github.com/yourusername/vaultpipe/internal/runbook"
	"github.com/yourusername/vaultpipe/internal/vault"
)

// stubClient satisfies the interface used by runbook via vault.Client.
// We rely on the real vault.Client being wrappable; here we use a test double
// by embedding a fake through build-tag-free interface extraction.

type fakeClient struct {
	pingErr  error
	readErr  error
	readData map[string]interface{}
}

func (f *fakeClient) Ping(_ context.Context) error { return f.pingErr }
func (f *fakeClient) Read(_ context.Context, _ string) (map[string]interface{}, error) {
	if f.readErr != nil {
		return nil, f.readErr
	}
	return f.readData, nil
}

func makeConfig(paths ...string) *config.Config {
	secrets := make([]config.SecretRef, len(paths))
	for i, p := range paths {
		secrets[i] = config.SecretRef{Path: p, Env: "KEY"}
	}
	return &config.Config{Secrets: secrets}
}

func TestRun_AllPassing(t *testing.T) {
	client := vault.NewClientFromFake(&fakeClient{readData: map[string]interface{}{"value": "s"}})
	r := runbook.New(client, makeConfig("secret/data/app"))
	results, err := r.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, res := range results {
		if !res.Passed {
			t.Errorf("check %q should pass, got: %s", res.Name, res.Message)
		}
	}
}

func TestRun_VaultUnreachable(t *testing.T) {
	client := vault.NewClientFromFake(&fakeClient{pingErr: errors.New("connection refused")})
	r := runbook.New(client, makeConfig())
	results, err := r.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Passed {
		t.Error("vault:reachable should fail when ping errors")
	}
}

func TestRun_SecretUnreadable(t *testing.T) {
	client := vault.NewClientFromFake(&fakeClient{readErr: errors.New("403 forbidden")})
	r := runbook.New(client, makeConfig("secret/data/missing"))
	results, _ := r.Run(context.Background())
	if results[1].Passed {
		t.Error("secret check should fail when read errors")
	}
}

func TestSummary_AllPass(t *testing.T) {
	results := []runbook.CheckResult{
		{Name: "vault:reachable", Passed: true},
		{Name: "secret:foo", Passed: true},
	}
	msg, ok := runbook.Summary(results)
	if !ok {
		t.Errorf("expected ok, got message: %s", msg)
	}
}

func TestSummary_SomeFail(t *testing.T) {
	results := []runbook.CheckResult{
		{Name: "vault:reachable", Passed: true},
		{Name: "secret:bar", Passed: false, Message: "403"},
	}
	msg, ok := runbook.Summary(results)
	if ok {
		t.Error("expected failure summary")
	}
	if msg == "" {
		t.Error("expected non-empty summary message")
	}
}
