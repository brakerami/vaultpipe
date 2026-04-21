package circuit_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultpipe/internal/circuit"
)

func TestNew_InvalidThreshold(t *testing.T) {
	_, err := circuit.New(0, time.Second)
	if err == nil {
		t.Fatal("expected error for zero threshold")
	}
}

func TestNew_InvalidResetTimeout(t *testing.T) {
	_, err := circuit.New(3, 0)
	if err == nil {
		t.Fatal("expected error for zero resetTimeout")
	}
}

func TestBreaker_InitiallyAllows(t *testing.T) {
	b, _ := circuit.New(3, time.Second)
	if !b.Allow() {
		t.Fatal("expected new breaker to allow calls")
	}
}

func TestBreaker_OpensAfterThreshold(t *testing.T) {
	b, _ := circuit.New(3, time.Second)
	b.RecordFailure()
	b.RecordFailure()
	if b.State() != circuit.StateClosed {
		t.Fatal("expected breaker to remain closed before threshold")
	}
	b.RecordFailure()
	if b.State() != circuit.StateOpen {
		t.Fatalf("expected StateOpen, got %v", b.State())
	}
	if b.Allow() {
		t.Fatal("expected open breaker to reject calls")
	}
}

func TestBreaker_SuccessResetsClosed(t *testing.T) {
	b, _ := circuit.New(2, time.Second)
	b.RecordFailure()
	b.RecordFailure()
	if b.State() != circuit.StateOpen {
		t.Fatal("expected open state")
	}
	// Simulate reset timeout by manipulating via half-open transition.
	// Use a very short timeout to allow the test to proceed.
	b2, _ := circuit.New(2, time.Millisecond)
	b2.RecordFailure()
	b2.RecordFailure()
	time.Sleep(5 * time.Millisecond)
	if !b2.Allow() {
		t.Fatal("expected half-open after timeout")
	}
	if b2.State() != circuit.StateHalfOpen {
		t.Fatalf("expected StateHalfOpen, got %v", b2.State())
	}
	b2.RecordSuccess()
	if b2.State() != circuit.StateClosed {
		t.Fatalf("expected StateClosed after success, got %v", b2.State())
	}
}

func TestBreaker_HalfOpen_FailureReopens(t *testing.T) {
	b, _ := circuit.New(1, time.Millisecond)
	b.RecordFailure()
	time.Sleep(5 * time.Millisecond)
	b.Allow() // transition to half-open
	b.RecordFailure()
	if b.State() != circuit.StateOpen {
		t.Fatalf("expected StateOpen after half-open failure, got %v", b.State())
	}
}
