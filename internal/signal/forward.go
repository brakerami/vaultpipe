// Package signal provides OS signal forwarding to child processes.
package signal

import (
	"os"
	"os/signal"
	"syscall"
)

// Forwarder listens for OS signals and forwards them to a target process.
type Forwarder struct {
	proc   *os.Process
	signals chan os.Signal
	done   chan struct{}
}

// New creates a Forwarder that will forward signals to the given process.
func New(proc *os.Process) *Forwarder {
	return &Forwarder{
		proc:    proc,
		signals: make(chan os.Signal, 8),
		done:    make(chan struct{}),
	}
}

// Start begins forwarding signals. Call Stop to clean up.
func (f *Forwarder) Start() {
	signal.Notify(f.signals,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
	)
	go f.loop()
}

// Stop stops forwarding signals and cleans up resources.
func (f *Forwarder) Stop() {
	signal.Stop(f.signals)
	close(f.done)
}

func (f *Forwarder) loop() {
	for {
		select {
		case sig := <-f.signals:
			if f.proc != nil {
				_ = f.proc.Signal(sig)
			}
		case <-f.done:
			return
		}
	}
}
