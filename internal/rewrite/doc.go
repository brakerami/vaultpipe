// Package rewrite implements value transformation rules for secrets.
//
// Rules are applied in declaration order before a secret value is injected
// into the process environment. Supported transformations are:
//
//   - upper   – convert to upper-case
//   - lower   – convert to lower-case
//   - trim    – strip leading/trailing whitespace
//   - prefix  – prepend a fixed string (requires Arg)
//   - suffix  – append a fixed string (requires Arg)
//
// Example usage:
//
//	rw, err := rewrite.New([]rewrite.Rule{
//		{Kind: "trim"},
//		{Kind: "prefix", Arg: "prod_"},
//	})
//	if err != nil { ... }
//	transformed := rw.Apply(rawSecretValue)
package rewrite
