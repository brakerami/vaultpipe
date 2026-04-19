package audit

import (
	"encoding/json"
	"io"
	"os"
	"time"
)

// Entry represents a single audit log event.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Event     string    `json:"event"`
	Path      string    `json:"path,omitempty"`
	EnvKey    string    `json:"env_key,omitempty"`
	Success   bool      `json:"success"`
	Error     string    `json:"error,omitempty"`
}

// Logger writes structured audit entries to a writer.
type Logger struct {
	out io.Writer
}

// NewLogger creates a Logger writing to w. If w is nil, os.Stderr is used.
func NewLogger(w io.Writer) *Logger {
	if w == nil {
		w = os.Stderr
	}
	return &Logger{out: w}
}

// Log writes a single audit entry as a JSON line.
func (l *Logger) Log(e Entry) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	_, err = l.out.Write(data)
	return err
}

// SecretFetched logs a secret resolution attempt.
func (l *Logger) SecretFetched(path, envKey string, err error) {
	e := Entry{
		Event:   "secret_fetched",
		Path:    path,
		EnvKey:  envKey,
		Success: err == nil,
	}
	if err != nil {
		e.Error = err.Error()
	}
	_ = l.Log(e)
}

// ProcessStarted logs that the child process is being launched.
func (l *Logger) ProcessStarted(cmd string) {
	_ = l.Log(Entry{
		Event:   "process_started",
		Path:    cmd,
		Success: true,
	})
}
