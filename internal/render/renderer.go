// Package render provides template rendering for environment config values.
// It supports Go text/template syntax with access to the full resolved chain.
package render

import (
	"bytes"
	"fmt"
	"text/template"
)

// Resolver is satisfied by any type that can look up a key's string value.
type Resolver interface {
	Get(key string) (string, bool)
}

// Renderer expands Go text/template expressions against a Resolver.
type Renderer struct {
	resolver Resolver
}

// NewRenderer returns a Renderer backed by the given Resolver.
func NewRenderer(r Resolver) *Renderer {
	return &Renderer{resolver: r}
}

// Render parses src as a Go text/template and executes it with a FuncMap
// that exposes an "env" function for looking up keys from the Resolver.
// Returns an error if parsing or execution fails, or if a required key is
// absent and no default was provided via the template itself.
func (r *Renderer) Render(src string) (string, error) {
	funcMap := template.FuncMap{
		"env": func(key string) (string, error) {
			v, ok := r.resolver.Get(key)
			if !ok {
				return "", fmt.Errorf("render: key %q not found", key)
			}
			return v, nil
		},
		"envOr": func(key, fallback string) string {
			v, ok := r.resolver.Get(key)
			if !ok {
				return fallback
			}
			return v
		},
	}

	tmpl, err := template.New("envchain").Funcs(funcMap).Parse(src)
	if err != nil {
		return "", fmt.Errorf("render: parse error: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, nil); err != nil {
		return "", fmt.Errorf("render: execute error: %w", err)
	}
	return buf.String(), nil
}
