// Package health provides a simple Vault connectivity check used by
// runbook diagnostics and startup validation.
package health

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Status represents the result of a Vault health check.
type Status struct {
	Reachable   bool
	Sealed      bool
	Standby     bool
	Version     string
	Latency     time.Duration
}

// Checker performs health checks against a Vault instance.
type Checker struct {
	baseURL string
	client  *http.Client
}

// New returns a Checker targeting the given Vault base URL.
func New(baseURL string, timeout time.Duration) *Checker {
	return &Checker{
		baseURL: baseURL,
		client: &http.Client{Timeout: timeout},
	}
}

// Check calls the Vault /v1/sys/health endpoint and returns a Status.
// A non-200 response is still parsed — Vault uses status codes to
// communicate sealed / standby state, so we treat them as structured
// information rather than errors.
func (c *Checker) Check(ctx context.Context) (Status, error) {
	url := c.baseURL + "/v1/sys/health?standbyok=true&sealedok=true"

	start := time.Now()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Status{}, fmt.Errorf("health: build request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return Status{Reachable: false}, fmt.Errorf("health: vault unreachable: %w", err)
	}
	defer resp.Body.Close()

	latency := time.Since(start)

	// Vault health status codes:
	//   200 — initialised, unsealed, active
	//   429 — unsealed, standby
	//   501 — not initialised
	//   503 — sealed
	var s Status
	s.Reachable = true
	s.Latency = latency
	s.Version = resp.Header.Get("X-Vault-Version")

	switch resp.StatusCode {
	case http.StatusOK:
		// active, healthy
	case http.StatusTooManyRequests:
		s.Standby = true
	case http.StatusServiceUnavailable:
		s.Sealed = true
	default:
		return s, fmt.Errorf("health: unexpected status %d", resp.StatusCode)
	}

	return s, nil
}
