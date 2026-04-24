package batch_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/yourusername/vaultpipe/internal/batch"
)

func TestCollect_NoErrors_ReturnsNil(t *testing.T) {
	results := []batch.Result{
		{Key: "A", Value: "v1"},
		{Key: "B", Value: "v2"},
	}
	if err := batch.Collect(results); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestCollect_WithErrors_ReturnsMultiError(t *testing.T) {
	results := []batch.Result{
		{Key: "A", Value: "v1"},
		{Key: "B", Err: errors.New("not found")},
		{Key: "C", Err: errors.New("forbidden")},
	}
	err := batch.Collect(results)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !batch.IsPartialFailure(err) {
		t.Fatalf("expected MultiError, got %T", err)
	}
	if !strings.Contains(err.Error(), "2 error(s)") {
		t.Errorf("unexpected message: %s", err.Error())
	}
}

func TestMultiError_Unwrap_ContainsWrappedErrors(t *testing.T) {
	sentinel := errors.New("sentinel")
	results := []batch.Result{
		{Key: "X", Err: sentinel},
	}
	err := batch.Collect(results)
	if !errors.Is(err, sentinel) {
		t.Errorf("expected errors.Is to find sentinel via Unwrap")
	}
}

func TestIsPartialFailure_NonMultiError(t *testing.T) {
	if batch.IsPartialFailure(errors.New("plain")) {
		t.Error("plain error should not be a partial failure")
	}
}

func TestIsPartialFailure_Nil(t *testing.T) {
	if batch.IsPartialFailure(nil) {
		t.Error("nil should not be a partial failure")
	}
}
