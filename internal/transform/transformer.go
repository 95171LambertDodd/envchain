// Package transform provides key/value transformation functions
// that can be applied to environment config layers, such as
// uppercasing keys, trimming whitespace, or prefixing values.
package transform

import (
	"fmt"
	"strings"
)

// TransformFunc is a function that transforms a key-value pair.
// It returns the transformed key and value, or an error.
type TransformFunc func(key, value string) (string, string, error)

// Transformer applies a chain of TransformFuncs to a map of entries.
type Transformer struct {
	fns []TransformFunc
}

// NewTransformer creates a Transformer with the given transform functions.
func NewTransformer(fns ...TransformFunc) (*Transformer, error) {
	if len(fns) == 0 {
		return nil, fmt.Errorf("transformer: at least one TransformFunc required")
	}
	return &Transformer{fns: fns}, nil
}

// Apply runs all transform functions over each entry in src and returns
// a new map with the transformed key-value pairs.
func (t *Transformer) Apply(src map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(src))
	for k, v := range src {
		curKey, curVal := k, v
		for _, fn := range t.fns {
			newKey, newVal, err := fn(curKey, curVal)
			if err != nil {
				return nil, fmt.Errorf("transformer: error on key %q: %w", curKey, err)
			}
			curKey, curVal = newKey, newVal
		}
		out[curKey] = curVal
	}
	return out, nil
}

// UppercaseKeys returns a TransformFunc that uppercases all keys.
func UppercaseKeys() TransformFunc {
	return func(key, value string) (string, string, error) {
		return strings.ToUpper(key), value, nil
	}
}

// TrimSpace returns a TransformFunc that trims whitespace from values.
func TrimSpace() TransformFunc {
	return func(key, value string) (string, string, error) {
		return key, strings.TrimSpace(value), nil
	}
}

// PrefixKeys returns a TransformFunc that prepends prefix to every key.
func PrefixKeys(prefix string) TransformFunc {
	return func(key, value string) (string, string, error) {
		return prefix + key, value, nil
	}
}

// StripPrefix returns a TransformFunc that removes a prefix from keys
// that have it, leaving other keys unchanged.
func StripPrefix(prefix string) TransformFunc {
	return func(key, value string) (string, string, error) {
		return strings.TrimPrefix(key, prefix), value, nil
	}
}
