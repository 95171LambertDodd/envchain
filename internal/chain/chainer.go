// Package chain provides utilities for chaining multiple key-value sources
// with priority ordering and optional fallback behaviour.
package chain

import (
	"errors"
	"fmt"
)

// Source is the minimal interface a key-value store must satisfy.
type Source interface {
	Keys() []string
	Get(key string) (string, bool)
}

// Chainer resolves keys across an ordered list of Sources.
// The first source that contains the key wins.
type Chainer struct {
	sources []Source
	labels  []string
}

// New creates a Chainer from the provided sources.
// At least one source must be supplied.
func New(sources ...Source) (*Chainer, error) {
	if len(sources) == 0 {
		return nil, errors.New("chainer: at least one source is required")
	}
	return &Chainer{sources: sources, labels: make([]string, len(sources))}, nil
}

// WithLabel attaches a human-readable label to the source at position idx.
// Labels are surfaced in error messages and debug output.
func (c *Chainer) WithLabel(idx int, label string) error {
	if idx < 0 || idx >= len(c.sources) {
		return fmt.Errorf("chainer: index %d out of range (have %d sources)", idx, len(c.sources))
	}
	c.labels[idx] = label
	return nil
}

// Get returns the value for key from the highest-priority source that has it.
// The second return value is false when no source contains the key.
func (c *Chainer) Get(key string) (string, bool) {
	for _, s := range c.sources {
		if v, ok := s.Get(key); ok {
			return v, true
		}
	}
	return "", false
}

// Keys returns the deduplicated union of all keys across every source.
func (c *Chainer) Keys() []string {
	seen := make(map[string]struct{})
	var out []string
	for _, s := range c.sources {
		for _, k := range s.Keys() {
			if _, exists := seen[k]; !exists {
				seen[k] = struct{}{}
				out = append(out, k)
			}
		}
	}
	return out
}

// Origin returns the label of the source that owns key, or an empty string
// when no source contains the key.
func (c *Chainer) Origin(key string) string {
	for i, s := range c.sources {
		if _, ok := s.Get(key); ok {
			if c.labels[i] != "" {
				return c.labels[i]
			}
			return fmt.Sprintf("source[%d]", i)
		}
	}
	return ""
}
