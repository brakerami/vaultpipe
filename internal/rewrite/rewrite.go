// Package rewrite provides value transformation rules applied to secrets
// before they are injected into the process environment.
package rewrite

import (
	"fmt"
	"strings"
)

// Rule describes a single transformation to apply to a secret value.
type Rule struct {
	// Kind is the transformation type: "upper", "lower", "trim", "prefix", "suffix".
	Kind string
	// Arg is an optional argument (e.g. the prefix/suffix string).
	Arg string
}

// Rewriter applies an ordered list of Rules to secret values.
type Rewriter struct {
	rules []Rule
}

// New returns a Rewriter that applies rules in order.
func New(rules []Rule) (*Rewriter, error) {
	for i, r := range rules {
		switch r.Kind {
		case "upper", "lower", "trim":
			// no arg required
		case "prefix", "suffix":
			if r.Arg == "" {
				return nil, fmt.Errorf("rewrite rule %d (%q) requires a non-empty arg", i, r.Kind)
			}
		default:
			return nil, fmt.Errorf("rewrite rule %d: unknown kind %q", i, r.Kind)
		}
	}
	return &Rewriter{rules: rules}, nil
}

// Apply runs all rules against value and returns the transformed result.
func (r *Rewriter) Apply(value string) string {
	for _, rule := range r.rules {
		switch rule.Kind {
		case "upper":
			value = strings.ToUpper(value)
		case "lower":
			value = strings.ToLower(value)
		case "trim":
			value = strings.TrimSpace(value)
		case "prefix":
			value = rule.Arg + value
		case "suffix":
			value = value + rule.Arg
		}
	}
	return value
}

// Len returns the number of rules registered.
func (r *Rewriter) Len() int { return len(r.rules) }
