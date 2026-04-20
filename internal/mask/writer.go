package mask

import "io"

// Writer wraps an io.Writer and masks secrets before writing.
type Writer struct {
	underlying io.Writer
	masker     *Masker
}

// NewWriter returns a Writer that redacts secrets from all writes.
func NewWriter(w io.Writer, m *Masker) *Writer {
	return &Writer{underlying: w, masker: m}
}

// Write masks p and writes the result to the underlying writer.
// The original length of p is returned so callers are not confused by
// the potentially shorter masked output.
func (w *Writer) Write(p []byte) (n int, err error) {
	masked := w.masker.Mask(string(p))
	_, err = w.underlying.Write([]byte(masked))
	if err != nil {
		return 0, err
	}
	// Return original length so callers are not confused.
	return len(p), nil
}

// Unwrap returns the underlying io.Writer. This is useful for callers that
// need to access the original writer after masking is no longer required,
// and follows the convention used by other wrapper types in the standard
// library (e.g. bufio.Writer).
func (w *Writer) Unwrap() io.Writer {
	return w.underlying
}
