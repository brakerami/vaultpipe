package signal_test

import (
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/signal"
)

func TestNew_ReturnsForwarder(t *testing.T) {
	f := signal.New(nil)
	if f == nil {
		t.Fatal("expected non-nil Forwarder")
	}
}

func TestForwarder_StartStop(t *testing.T) {
	f := signal.New(nil)
	f.Start()
	// Should not block or panic
	f.Stop()
}

func TestForwarder_ForwardsSignalToProcess(t *testing.T) {
	// Start a real child process that sleeps and records signals.
	cmd := exec.Command("sleep", "10")
	if err := cmd.Start(); err != nil {
		t.Skipf("could not start sleep: %v", err)
	}
	defer func() { _ = cmd.Process.Kill() }()

	f := signal.New(cmd.Process)
	f.Start()
	defer f.Stop()

	// Send SIGTERM to ourselves; the forwarder should relay it to the child.
	self, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatalf("could not find self: %v", err)
	}
	_ = self.Signal(syscall.SIGTERM)

	// Give the goroutine time to forward.
	time.Sleep(100 * time.Millisecond)

	// Child should have exited due to SIGTERM.
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	select {
	case <-done:
		// success — child received the signal
	case <-time.After(2 * time.Second):
		t.Error("child process did not exit after forwarded SIGTERM")
	}
}
