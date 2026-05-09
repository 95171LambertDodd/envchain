// Package trim provides utilities for trimming and normalising keys and values
// within an envchain layer before further processing.
package trim

import (
	"fmt"
	"strings"
)

// Source is the minimal read interface required by Trimmer.
type Source interface {
	Keys() []string
	Get(key string) (string, bool)
}

// TrimMode controls which parts of an entry are trimmed.
type TrimMode int

const (
	TrimBoth   TrimMode = iota // trim leading and trailing whitespace from values
	TrimKeys                   // trim whitespace from keys (produces new key names)
	TrimAll                    // trim both keys and values
)

// Trimmer applies whitespace trimming to keys and/or values from a Source.
type Trimmer struct {
	src  Source
	mode TrimMode
}

// NewTrimmer constructs a Trimmer for the given source and mode.
// Returns an error if src is nil.
func NewTrimmer(src Source, mode TrimMode) (*Trimmer, error) {
	if src == nil {
		return nil, fmt.Errorf("trim: source must not be nil")
	}
	return &Trimmer{src: src, mode: mode}, nil
}

// Apply iterates over all keys in the source, trims according to the mode,
// and returns the resulting map. Duplicate trimmed keys are last-write-wins.
func (t *Trimmer) Apply() map[string]string {
	out := make(map[string]string)
	for _, k := range t.src.Keys() {
		v, ok := t.src.Get(k)
		if !ok {
			continue
		}
		outKey := k
		outVal := v
		if t.mode == TrimKeys || t.mode == TrimAll {
			outKey = strings.TrimSpace(k)
		}
		if t.mode == TrimBoth || t.mode == TrimAll {
			outVal = strings.TrimSpace(v)
		}
		if outKey == "" {
			continue
		}
		out[outKey] = outVal
	}
	return out
}

// ApplyToLayer trims all values (TrimBoth) and writes results into dest via Set.
// dest must expose a Set(key, value string) error method.
func ApplyToLayer(src Source, dest interface {
	Set(key, value string) error
}) error {
	t, err := NewTrimmer(src, TrimBoth)
	if err != nil {
		return err
	}
	for k, v := range t.Apply() {
		if err := dest.Set(k, v); err != nil {
			return fmt.Errorf("trim: setting key %q: %w", k, err)
		}
	}
	return nil
}
