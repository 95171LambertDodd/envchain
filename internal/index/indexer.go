// Package index provides key-based indexing over environment layers,
// allowing fast lookup, listing, and reverse-lookup of keys by value.
package index

import (
	"errors"
	"fmt"
	"sort"
)

// Source is the minimal interface required to build an index.
type Source interface {
	Keys() []string
	Get(key string) (string, bool)
}

// Indexer maintains a forward (key→value) and reverse (value→keys) index
// over a Source snapshot taken at construction time.
type Indexer struct {
	forward map[string]string
	reverse map[string][]string
}

// NewIndexer builds an Indexer from the given Source.
// Returns an error if source is nil.
func NewIndexer(src Source) (*Indexer, error) {
	if src == nil {
		return nil, errors.New("index: source must not be nil")
	}

	fwd := make(map[string]string)
	rev := make(map[string][]string)

	for _, k := range src.Keys() {
		v, ok := src.Get(k)
		if !ok {
			continue
		}
		fwd[k] = v
		rev[v] = append(rev[v], k)
	}

	// Sort reverse-lookup slices for deterministic output.
	for v := range rev {
		sort.Strings(rev[v])
	}

	return &Indexer{forward: fwd, reverse: rev}, nil
}

// Get returns the value for a key and whether it was found.
func (idx *Indexer) Get(key string) (string, bool) {
	v, ok := idx.forward[key]
	return v, ok
}

// KeysForValue returns all keys that map to the given value, sorted.
// Returns an empty slice when no keys share that value.
func (idx *Indexer) KeysForValue(value string) []string {
	keys, ok := idx.reverse[value]
	if !ok {
		return []string{}
	}
	out := make([]string, len(keys))
	copy(out, keys)
	return out
}

// Keys returns all indexed keys in sorted order.
func (idx *Indexer) Keys() []string {
	keys := make([]string, 0, len(idx.forward))
	for k := range idx.forward {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// HasDuplicateValues reports whether any two keys share the same value
// and returns one example pair for diagnostics.
func (idx *Indexer) HasDuplicateValues() (bool, string) {
	for v, keys := range idx.reverse {
		if len(keys) > 1 {
			return true, fmt.Sprintf("value %q shared by keys: %v", v, keys)
		}
	}
	return false, ""
}
