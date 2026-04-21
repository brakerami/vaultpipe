package pipeline_test

import (
	"context"
	"errors"
	"testing"

	"github.com/your-org/vaultpipe/internal/environ"
	"github.com/your-org/vaultpipe/internal/pipeline"
	"github.com/your-org/vaultpipe/internal/redact"
	"github.com/your-org/vaultpipe/internal/resolver"
)

// --- stubs ---

type stubResolver struct {
	vals map[string]string
	err  error
}

func (s *stubResolver) Resolve(_ context.Context, ref string) (string, error) {
	if s.err != nil {
		return "", s.err
	}
	return s.vals[ref], nil
}

type stubRunner struct {
	gotEnv []string
	err    error
}

func (s *stubRunner) Run(_ context.Context, env []string) error {
	s.gotEnv = env
	return s.err
}

// --- tests ---

func TestPipeline_Run_InjectsSecrets(t *testing.T) {
	res := &stubResolver{vals: map[string]string{"secret/db#pass": "hunter2"}}
	run := &stubRunner{}
	refs := []resolver.Ref{{Env: "DB_PASS", Path: "secret/db#pass"}}

	p := pipeline.New(refs, res, run, redact.New())
	base := environ.FromMap(map[string]string{"HOME": "/root"})

	if err := p.Run(context.Background(), base); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	envMap := toMap(run.gotEnv)
	if envMap["DB_PASS"] != "hunter2" {
		t.Errorf("expected DB_PASS=hunter2, got %q", envMap["DB_PASS"])
	}
	if envMap["HOME"] != "/root" {
		t.Errorf("base env not preserved, HOME=%q", envMap["HOME"])
	}
}

func TestPipeline_Run_ResolveError(t *testing.T) {
	res := &stubResolver{err: errors.New("vault down")}
	run := &stubRunner{}
	refs := []resolver.Ref{{Env: "TOKEN", Path: "secret/token"}}

	p := pipeline.New(refs, res, run, redact.New())
	err := p.Run(context.Background(), environ.FromMap(nil))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestPipeline_Run_RunnerError(t *testing.T) {
	res := &stubResolver{vals: map[string]string{"secret/x": "val"}}
	run := &stubRunner{err: errors.New("exec failed")}
	refs := []resolver.Ref{{Env: "X", Path: "secret/x"}}

	p := pipeline.New(refs, res, run, redact.New())
	err := p.Run(context.Background(), environ.FromMap(nil))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func toMap(pairs []string) map[string]string {
	m := make(map[string]string, len(pairs))
	for _, p := range pairs {
		for i := 0; i < len(p); i++ {
			if p[i] == '=' {
				m[p[:i]] = p[i+1:]
				break
			}
		}
	}
	return m
}
