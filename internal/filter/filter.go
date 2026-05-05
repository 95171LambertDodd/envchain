// Package filter provides key-based filtering of environment layers,
// supporting prefix matching, suffix matching, and exact key sets.
package filter

import (
	"strings"
)

// Mode controls how Filter matches keys.
type Mode int

const (
	ModePrefix Mode = iota
	ModeSuffix
	ModeExact
)

// Filter selects a subset of key-value pairs from a map.
type Filter struct {
	mode    Mode
	patterns []string
}

// NewFilter creates a Filter with the given mode and patterns.
// Returns an error if no patterns are provided.
func NewFilter(mode Mode, patterns ...string) (*Filter, error) {
	if len(patterns) == 0 {
		return nil, ErrNoPatternsProvided
	}
	return &Filter{mode: mode, patterns: patterns}, nil
}

// Apply returns a new map containing only entries whose keys match
// at least one of the filter's patterns.
func (f *Filter) Apply(env map[string]string) map[string]string {
	out := make(map[string]string)
	for k, v := range env {
		if f.matches(k) {
			out[k] = v
		}
	}
	return out
}

// Exclude returns a new map with matching keys removed.
func (f *Filter) Exclude(env map[string]string) map[string]string {
	out := make(map[string]string)
	for k, v := range env {
		if !f.matches(k) {
			out[k] = v
		}
	}
	return out
}

func (f *Filter) matches(key string) bool {
	for _, p := range f.patterns {
		switch f.mode {
		case ModePrefix:
			if strings.HasPrefix(key, p) {
				return true
			}
		case ModeSuffix:
			if strings.HasSuffix(key, p) {
				return true
			}
		case ModeExact:
			if key == p {
				return true
			}
		}
	}
	return false
}
