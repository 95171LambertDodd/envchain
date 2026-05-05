// Package mask provides utilities for redacting sensitive environment
// variable values before they are logged, exported, or displayed.
package mask

import (
	"strings"
)

const redacted = "[REDACTED]"

// Masker holds a set of key patterns considered sensitive and replaces
// their values with a fixed redaction string.
type Masker struct {
	patterns []string
}

// NewMasker returns a Masker that treats any key whose lower-cased name
// contains one of the provided patterns as sensitive.
func NewMasker(patterns ...string) *Masker {
	norm := make([]string, len(patterns))
	for i, p := range patterns {
		norm[i] = strings.ToLower(p)
	}
	return &Masker{patterns: norm}
}

// IsSensitive reports whether the given key should be masked.
func (m *Masker) IsSensitive(key string) bool {
	lower := strings.ToLower(key)
	for _, p := range m.patterns {
		if strings.Contains(lower, p) {
			return true
		}
	}
	return false
}

// MaskValue returns the redaction placeholder if the key is sensitive,
// otherwise it returns the original value unchanged.
func (m *Masker) MaskValue(key, value string) string {
	if m.IsSensitive(key) {
		return redacted
	}
	return value
}

// MaskMap returns a new map with sensitive values replaced.
// The original map is never modified.
func (m *Masker) MaskMap(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = m.MaskValue(k, v)
	}
	return out
}
