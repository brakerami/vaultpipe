package health_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/health"
)

func newTestChecker(t *testing.T, statusCode int, vaultVersion string) (*health.Checker, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if vaultVersion != "" {
			w.Header().Set("X-Vault-Version", vaultVersion)
		}
		w.WriteHeader(statusCode)
	}))
	t.Cleanup(srv.Close)
	checker := health.New(srv.URL, 5*time.Second)
	return checker, srv
}

func TestCheck_ActiveVault(t *testing.T) {
	checker, _ := newTestChecker(t, http.StatusOK, "1.15.0")

	status, err := checker.Check(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !status.Reachable {
		t.Error("expected Reachable=true")
	}
	if status.Sealed {
		t.Error("expected Sealed=false")
	}
	if status.Standby {
		t.Error("expected Standby=false")
	}
	if status.Version != "1.15.0" {
		t.Errorf("expected version 1.15.0, got %q", status.Version)
	}
}

func TestCheck_SealedVault(t *testing.T) {
	checker, _ := newTestChecker(t, http.StatusServiceUnavailable, "")

	status, err := checker.Check(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !status.Reachable {
		t.Error("expected Reachable=true even when sealed")
	}
	if !status.Sealed {
		t.Error("expected Sealed=true")
	}
}

func TestCheck_StandbyVault(t *testing.T) {
	checker, _ := newTestChecker(t, http.StatusTooManyRequests, "")

	status, err := checker.Check(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !status.Standby {
		t.Error("expected Standby=true")
	}
}

func TestCheck_Unreachable(t *testing.T) {
	checker := health.New("http://127.0.0.1:19999", 200*time.Millisecond)

	status, err := checker.Check(context.Background())
	if err == nil {
		t.Fatal("expected error for unreachable vault")
	}
	if status.Reachable {
		t.Error("expected Reachable=false")
	}
}

func TestCheck_UnexpectedStatusCode(t *testing.T) {
	checker, _ := newTestChecker(t, http.StatusTeapot, "")

	_, err := checker.Check(context.Background())
	if err == nil {
		t.Fatal("expected error for unexpected status code")
	}
}

func TestCheck_LatencyPopulated(t *testing.T) {
	checker, _ := newTestChecker(t, http.StatusOK, "")

	status, err := checker.Check(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status.Latency <= 0 {
		t.Error("expected positive latency")
	}
}
