// Package coerce provides type coercion utilities for environment config values.
// It converts string values from a config source into typed Go values with
// optional default fallbacks.
package coerce

import (
	"fmt"
	"strconv"
	"strings"
)

// Source is the interface for retrieving raw string config values.
type Source interface {
	Get(key string) (string, bool)
}

// Coercer wraps a Source and provides typed accessors.
type Coercer struct {
	src Source
}

// NewCoercer returns a Coercer backed by src.
// Returns an error if src is nil.
func NewCoercer(src Source) (*Coercer, error) {
	if src == nil {
		return nil, fmt.Errorf("coerce: source must not be nil")
	}
	return &Coercer{src: src}, nil
}

// String returns the raw string value for key, or fallback if not found.
func (c *Coercer) String(key, fallback string) string {
	if v, ok := c.src.Get(key); ok {
		return v
	}
	return fallback
}

// Int parses the value for key as an integer.
// Returns fallback if the key is absent. Returns an error if parsing fails.
func (c *Coercer) Int(key string, fallback int) (int, error) {
	v, ok := c.src.Get(key)
	if !ok {
		return fallback, nil
	}
	n, err := strconv.Atoi(strings.TrimSpace(v))
	if err != nil {
		return 0, fmt.Errorf("coerce: key %q: cannot parse %q as int", key, v)
	}
	return n, nil
}

// Bool parses the value for key as a boolean (true/false/1/0/yes/no).
// Returns fallback if the key is absent. Returns an error if parsing fails.
func (c *Coercer) Bool(key string, fallback bool) (bool, error) {
	v, ok := c.src.Get(key)
	if !ok {
		return fallback, nil
	}
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "true", "1", "yes":
		return true, nil
	case "false", "0", "no":
		return false, nil
	default:
		return false, fmt.Errorf("coerce: key %q: cannot parse %q as bool", key, v)
	}
}

// Float parses the value for key as a float64.
// Returns fallback if the key is absent. Returns an error if parsing fails.
func (c *Coercer) Float(key string, fallback float64) (float64, error) {
	v, ok := c.src.Get(key)
	if !ok {
		return fallback, nil
	}
	f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
	if err != nil {
		return 0, fmt.Errorf("coerce: key %q: cannot parse %q as float", key, v)
	}
	return f, nil
}
