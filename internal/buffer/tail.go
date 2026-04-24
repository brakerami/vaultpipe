package buffer

import (
	"fmt"
	"io"
	"strings"
)

// Tail writes the most recent n entries from r to w in a human-readable
// format. If n <= 0 all entries are written. Returns the number of lines
// written and any error from w.
func Tail(w io.Writer, r *Ring, n int) (int, error) {
	entries := r.Snapshot()
	if n > 0 && n < len(entries) {
		entries = entries[len(entries)-n:]
	}
	written := 0
	for _, e := range entries {
		line := fmt.Sprintf("%s  %s\n", e.At.Format("15:04:05.000"), e.Message)
		if _, err := io.WriteString(w, line); err != nil {
			return written, err
		}
		written++
	}
	return written, nil
}

// TailString is a convenience wrapper that returns the tail as a string.
func TailString(r *Ring, n int) string {
	var sb strings.Builder
	_, _ = Tail(&sb, r, n)
	return sb.String()
}
