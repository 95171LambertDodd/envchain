// Package redact provides utilities for redacting sensitive environment
// variable values before logging, exporting, or displaying configuration.
package redact

import "strings"

// RedactMode controls how sensitive values are redacted.
type RedactMode int

const (
	// ModeStars replaces the value with a fixed string of asterisks.
	ModeStars RedactMode = iota
	// ModePartial reveals the first two characters and masks the rest.
	ModePartial
	// ModeHash replaces the value with a length-encoded placeholder.
	ModeHash
)

// Redactor redacts sensitive keys in a key/value map.
type Redactor struct {
	sensitiveSubstrings []string
	mode                RedactMode
}

// NewRedactor creates a Redactor that treats keys containing any of the
// provided substrings (case-insensitive) as sensitive.
func NewRedactor(sensitive []string, mode RedactMode) *Redactor {
	norm := make([]string, len(sensitive))
	for i, s := range sensitive {
		norm[i] = strings.ToLower(s)
	}
	return &Redactor{sensitiveSubstrings: norm, mode: mode}
}

// IsSensitive reports whether the given key should be redacted.
func (r *Redactor) IsSensitive(key string) bool {
	lower := strings.ToLower(key)
	for _, sub := range r.sensitiveSubstrings {
		if strings.Contains(lower, sub) {
			return true
		}
	}
	return false
}

// Redact returns the redacted form of value according to the configured mode.
func (r *Redactor) Redact(value string) string {
	switch r.mode {
	case ModePartial:
		if len(value) <= 2 {
			return "**"
		}
		return value[:2] + strings.Repeat("*", len(value)-2)
	case ModeHash:
		return "[redacted:" + strings.Repeat("*", min(len(value), 8)) + "]"
	default: // ModeStars
		return "********"
	}
}

// RedactMap returns a copy of m with sensitive values replaced.
func (r *Redactor) RedactMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		if r.IsSensitive(k) {
			out[k] = r.Redact(v)
		} else {
			out[k] = v
		}
	}
	return out
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
