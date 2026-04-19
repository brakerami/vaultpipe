package audit_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"vaultpipe/internal/audit"
)

func TestNewLogger_DefaultsToStderr(t *testing.T) {
	l := audit.NewLogger(nil)
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestLog_WritesJSON(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf)

	now := time.Now().UTC()
	err := l.Log(audit.Entry{
		Timestamp: now,
		Event:     "test_event",
		Success:   true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	line := strings.TrimSpace(buf.String())
	var got audit.Entry
	if err := json.Unmarshal([]byte(line), &got); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}
	if got.Event != "test_event" {
		t.Errorf("expected event test_event, got %s", got.Event)
	}
	if !got.Success {
		t.Error("expected success true")
	}
}

func TestSecretFetched_Success(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf)
	l.SecretFetched("secret/data/app", "DB_PASS", nil)

	var got audit.Entry
	json.Unmarshal(buf.Bytes(), &got)
	if got.Event != "secret_fetched" || !got.Success {
		t.Errorf("unexpected entry: %+v", got)
	}
}

func TestSecretFetched_Error(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf)
	l.SecretFetched("secret/data/app", "DB_PASS", errors.New("permission denied"))

	var got audit.Entry
	json.Unmarshal(buf.Bytes(), &got)
	if got.Success {
		t.Error("expected success false")
	}
	if got.Error != "permission denied" {
		t.Errorf("unexpected error field: %s", got.Error)
	}
}

func TestLog_SetsTimestampIfZero(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf)
	_ = l.Log(audit.Entry{Event: "no_ts"})

	var got audit.Entry
	json.Unmarshal(buf.Bytes(), &got)
	if got.Timestamp.IsZero() {
		t.Error("expected timestamp to be set automatically")
	}
}
