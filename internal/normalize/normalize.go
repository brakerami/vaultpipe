// Package normalize provides key normalization for environment variable names.
// It transforms arbitrary string keys into valid, consistent environment variable
// names by uppercasing, replacing illegal characters, and collapsing runs of
// separators.
package normalize

import (
	"strings"
	"unicode"
)

// Option controls how normalization is applied.
type Option func(*options)

type options struct {
	prefix    string
	separator rune
}

// WithPrefix prepends a fixed string (already normalized) to every key.
func WithPrefix(p string) Option {
	return func(o *options) { o.prefix = p }
}

// WithSeparator overrides the default underscore separator.
func WithSeparator(r rune) Option {
	return func(o *options) { o.separator = r }
}

// Key normalizes s into a valid environment variable name.
//
// Rules applied in order:
//  1. Convert to upper-case.
//  2. Replace any character that is not [A-Z0-9] with the separator.
//  3. Collapse consecutive separators into one.
//  4. Trim leading/trailing separators.
//  5. Prepend prefix when set.
func Key(s string, opts ...Option) string {
	o := &options{separator: '_'}
	for _, opt := range opts {
		opt(o)
	}

	s = strings.ToUpper(s)

	var b strings.Builder
	b.Grow(len(s))

	prevSep := true // treat start as separator to trim leading ones
	for _, r := range s {
		if unicode.IsUpper(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			prevSep = false
		} else {
			if !prevSep {
				b.WriteRune(o.separator)
			}
			prevSep = true
		}
	}

	result := strings.TrimRight(b.String(), string(o.separator))

	if o.prefix != "" {
		return o.prefix + string(o.separator) + result
	}
	return result
}

// Map normalizes every key in m, returning a new map.
// If two source keys collide after normalization the later value (iteration
// order is not guaranteed) wins.
func Map(m map[string]string, opts ...Option) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[Key(k, opts...)] = v
	}
	return out
}
