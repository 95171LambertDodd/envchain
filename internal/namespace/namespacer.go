// Package namespace provides key namespacing for environment config layers,
// allowing keys to be prefixed with a namespace and later stripped or queried
// by namespace.
package namespace

import (
	"errors"
	"fmt"
	"strings"
)

// Source is the minimal interface required to read environment entries.
type Source interface {
	Keys() []string
	Get(key string) (string, bool)
}

// Namespacer wraps a Source and applies a namespace prefix to all keys.
type Namespacer struct {
	source    Source
	namespace string
	sep       string
}

// NewNamespacer creates a Namespacer that prefixes keys with namespace+sep.
// namespace and sep must both be non-empty.
func NewNamespacer(source Source, namespace, sep string) (*Namespacer, error) {
	if source == nil {
		return nil, errors.New("namespace: source must not be nil")
	}
	if namespace == "" {
		return nil, errors.New("namespace: namespace must not be empty")
	}
	if sep == "" {
		return nil, errors.New("namespace: separator must not be empty")
	}
	return &Namespacer{source: source, namespace: namespace, sep: sep}, nil
}

// Prefix returns the full prefix string (namespace + separator).
func (n *Namespacer) Prefix() string {
	return n.namespace + n.sep
}

// QualifiedKey returns the namespaced form of key.
func (n *Namespacer) QualifiedKey(key string) string {
	return fmt.Sprintf("%s%s%s", n.namespace, n.sep, key)
}

// Strip removes the namespace prefix from key if present, returning the bare
// key and true. If the key does not belong to this namespace, it returns the
// original key and false.
func (n *Namespacer) Strip(key string) (string, bool) {
	prefix := n.Prefix()
	if strings.HasPrefix(key, prefix) {
		return strings.TrimPrefix(key, prefix), true
	}
	return key, false
}

// Keys returns all source keys qualified with the namespace prefix.
func (n *Namespacer) Keys() []string {
	raw := n.source.Keys()
	out := make([]string, len(raw))
	for i, k := range raw {
		out[i] = n.QualifiedKey(k)
	}
	return out
}

// Get retrieves a value by its qualified (namespaced) key. If the key does not
// belong to this namespace the lookup is skipped and false is returned.
func (n *Namespacer) Get(qualifiedKey string) (string, bool) {
	bare, ok := n.Strip(qualifiedKey)
	if !ok {
		return "", false
	}
	return n.source.Get(bare)
}

// Flatten returns a map of all qualified keys to their values.
func (n *Namespacer) Flatten() map[string]string {
	out := make(map[string]string)
	for _, k := range n.source.Keys() {
		if v, ok := n.source.Get(k); ok {
			out[n.QualifiedKey(k)] = v
		}
	}
	return out
}
