package redact

import "io"

// Writer wraps an io.Writer and scrubs secret values from every Write call.
type Writer struct {
	inner   io.Writer
	redactor *Redactor
}

// NewWriter returns a Writer that redacts secrets before forwarding to w.
func NewWriter(w io.Writer, r *Redactor) *Writer {
	return &Writer{inner: w, redactor: r}
}

// Write scrubs p and writes the result to the underlying writer.
func (w *Writer) Write(p []byte) (int, error) {
	clean := w.redactor.Scrub(string(p))
	_, err := w.inner.Write([]byte(clean))
	if err != nil {
		return 0, err
	}
	// Return original length so callers don't see a short-write error.
	return len(p), nil
}
