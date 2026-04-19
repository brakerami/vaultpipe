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
func (w *Writer) Write(p []byte) (n int, err error) {
	masked := w.masker.Mask(string(p))
	_, err = w.underlying.Write([]byte(masked))
	if err != nil {
		return 0, err
	}
	// Return original length so callers are not confused.
	return len(p), nil
}
