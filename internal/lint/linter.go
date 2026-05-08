// Package lint provides rule-based linting for environment config layers.
package lint

import (
	"fmt"
	"regexp"
	"strings"
)

// Severity indicates the importance of a lint finding.
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
	SeverityInfo    Severity = "info"
)

// Finding represents a single lint result.
type Finding struct {
	Key      string
	Rule     string
	Message  string
	Severity Severity
}

// Rule is a function that inspects a key/value pair and returns findings.
type Rule func(key, value string) []Finding

// Linter applies a set of rules to a map of environment entries.
type Linter struct {
	rules []Rule
}

// NewLinter creates a Linter with the provided rules.
// Returns an error if no rules are supplied.
func NewLinter(rules ...Rule) (*Linter, error) {
	if len(rules) == 0 {
		return nil, fmt.Errorf("lint: at least one rule is required")
	}
	return &Linter{rules: rules}, nil
}

// Lint runs all rules against every key/value pair in entries.
func (l *Linter) Lint(entries map[string]string) []Finding {
	var findings []Finding
	for k, v := range entries {
		for _, rule := range l.rules {
			findings = append(findings, rule(k, v)...)
		}
	}
	return findings
}

// NoEmptyValues flags keys whose value is an empty string.
func NoEmptyValues(key, value string) []Finding {
	if value == "" {
		return []Finding{{
			Key:      key,
			Rule:     "no-empty-values",
			Message:  fmt.Sprintf("key %q has an empty value", key),
			Severity: SeverityWarning,
		}}
	}
	return nil
}

// UppercaseKeys flags keys that contain lowercase letters.
func UppercaseKeys(key, _ string) []Finding {
	if key != strings.ToUpper(key) {
		return []Finding{{
			Key:      key,
			Rule:     "uppercase-keys",
			Message:  fmt.Sprintf("key %q should be uppercase", key),
			Severity: SeverityError,
		}}
	}
	return nil
}

// KeyPattern returns a Rule that enforces keys match the given regexp.
func KeyPattern(pattern string) (Rule, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("lint: invalid key pattern %q: %w", pattern, err)
	}
	return func(key, _ string) []Finding {
		if !re.MatchString(key) {
			return []Finding{{
				Key:      key,
				Rule:     "key-pattern",
				Message:  fmt.Sprintf("key %q does not match pattern %q", key, pattern),
				Severity: SeverityError,
			}}
		}
		return nil
	}, nil
}
