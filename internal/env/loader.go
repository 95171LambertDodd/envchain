// Package env provides utilities for loading environment variables
// into envchain layers, with support for prefix filtering and
// optional strict mode that errors on missing required keys.
package env

import (
	"fmt"
	"os"
	"strings"

	"github.com/envchain/envchain/internal/config"
)

// LoaderOption configures a Loader.
type LoaderOption func(*Loader)

// WithPrefix restricts loading to variables that start with the given prefix.
// The prefix is stripped from the key before storing in the layer.
func WithPrefix(prefix string) LoaderOption {
	return func(l *Loader) {
		l.prefix = prefix
	}
}

// WithStrict enables strict mode: Load will return an error if any key in
// required is absent from the environment.
func WithStrict(required ...string) LoaderOption {
	return func(l *Loader) {
		l.required = append(l.required, required...)
	}
}

// Loader reads OS environment variables into a config.Layer.
type Loader struct {
	prefix   string
	required []string
}

// NewLoader constructs a Loader with the supplied options.
func NewLoader(opts ...LoaderOption) *Loader {
	l := &Loader{}
	for _, o := range opts {
		o(l)
	}
	return l
}

// Load reads os.Environ, optionally filters by prefix, and populates a new
// Layer named layerName. It returns an error if strict mode is enabled and
// any required key is missing after filtering.
func (l *Loader) Load(layerName string) (*config.Layer, error) {
	if layerName == "" {
		return nil, fmt.Errorf("env: layer name must not be empty")
	}
	layer, err := config.NewLayer(layerName)
	if err != nil {
		return nil, fmt.Errorf("env: %w", err)
	}

	for _, entry := range os.Environ() {
		key, val, ok := strings.Cut(entry, "=")
		if !ok {
			continue
		}
		if l.prefix != "" {
			if !strings.HasPrefix(key, l.prefix) {
				continue
			}
			key = strings.TrimPrefix(key, l.prefix)
		}
		if setErr := layer.Set(key, val); setErr != nil {
			return nil, fmt.Errorf("env: set %q: %w", key, setErr)
		}
	}

	for _, req := range l.required {
		if _, err := layer.Get(req); err != nil {
			return nil, fmt.Errorf("env: required key %q not found in environment", req)
		}
	}

	return layer, nil
}
