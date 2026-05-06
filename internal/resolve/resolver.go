// Package resolve provides environment variable resolution across a chain
// of layers, supporting fallback values and strict missing-key behaviour.
package resolve

import (
	"errors"
	"fmt"
	"strings"
)

// ErrKeyNotFound is returned when a key cannot be resolved from any source.
var ErrKeyNotFound = errors.New("key not found")

// Source is the minimal interface required to look up a key.
type Source interface {
	Get(key string) (string, bool)
}

// Result holds the resolved value and the name of the source that provided it.
type Result struct {
	Key    string
	Value  string
	Source string
}

// Resolver resolves keys against an ordered list of named sources.
type Resolver struct {
	sources []namedSource
}

type namedSource struct {
	name string
	src  Source
}

// NewResolver creates a Resolver from an ordered slice of (name, Source) pairs.
// Sources are queried in order; the first match wins.
func NewResolver(sources ...struct {
	Name string
	Src  Source
}) *Resolver {
	r := &Resolver{}
	for _, s := range sources {
		r.sources = append(r.sources, namedSource{name: s.Name, src: s.Src})
	}
	return r
}

// Resolve returns the first value found for key across all sources.
// If no source contains the key, ErrKeyNotFound is returned.
func (r *Resolver) Resolve(key string) (Result, error) {
	for _, ns := range r.sources {
		if v, ok := ns.src.Get(key); ok {
			return Result{Key: key, Value: v, Source: ns.name}, nil
		}
	}
	return Result{}, fmt.Errorf("%w: %q", ErrKeyNotFound, key)
}

// ResolveAll resolves every key in keys and returns results or a combined error.
func (r *Resolver) ResolveAll(keys []string) ([]Result, error) {
	var results []Result
	var missing []string
	for _, k := range keys {
		res, err := r.Resolve(k)
		if err != nil {
			missing = append(missing, k)
			continue
		}
		results = append(results, res)
	}
	if len(missing) > 0 {
		return results, fmt.Errorf("%w: [%s]", ErrKeyNotFound, strings.Join(missing, ", "))
	}
	return results, nil
}
