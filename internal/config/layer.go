package config

import (
	"fmt"
	"os"
	"strings"
)

// Env represents a target environment tier.
type Env string

const (
	EnvDev     Env = "dev"
	EnvStaging Env = "staging"
	EnvProd    Env = "prod"
)

// Layer holds a named set of environment variables for a specific tier.
type Layer struct {
	Name   string
	Env    Env
	Values map[string]string
}

// NewLayer creates an empty Layer for the given environment.
func NewLayer(name string, env Env) *Layer {
	return &Layer{
		Name:   name,
		Env:    env,
		Values: make(map[string]string),
	}
}

// Set adds or updates a key-value pair in the layer.
func (l *Layer) Set(key, value string) error {
	key = strings.TrimSpace(key)
	if key == "" {
		return fmt.Errorf("config key must not be empty")
	}
	l.Values[key] = value
	return nil
}

// Get retrieves a value by key. Returns the value and whether it was found.
func (l *Layer) Get(key string) (string, bool) {
	v, ok := l.Values[key]
	return v, ok
}

// ApplyToEnv writes all key-value pairs in the layer to the current process
// environment. Existing variables are overwritten.
func (l *Layer) ApplyToEnv() error {
	for k, v := range l.Values {
		if err := os.Setenv(k, v); err != nil {
			return fmt.Errorf("failed to set env var %q: %w", k, err)
		}
	}
	return nil
}

// Keys returns a sorted list of all keys defined in the layer.
func (l *Layer) Keys() []string {
	keys := make([]string, 0, len(l.Values))
	for k := range l.Values {
		keys = append(keys, k)
	}
	return keys
}
