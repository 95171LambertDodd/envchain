// Package group provides key grouping functionality for envchain layers.
// Keys can be grouped by prefix, suffix, or a custom classifier function,
// enabling structured access to related configuration entries.
package group

import (
	"errors"
	"strings"
)

// Source is the interface satisfied by any key-value store.
type Source interface {
	Keys() []string
	Get(key string) (string, bool)
}

// GroupMode controls how keys are assigned to groups.
type GroupMode int

const (
	GroupByPrefix GroupMode = iota
	GroupBySuffix
	GroupByClassifier
)

// Grouper partitions keys from a Source into named buckets.
type Grouper struct {
	source     Source
	mode       GroupMode
	sep        string
	classifier func(key string) string
}

// NewGrouper constructs a Grouper for the given source.
// sep is the delimiter used when mode is GroupByPrefix or GroupBySuffix.
// classifier is required when mode is GroupByClassifier.
func NewGrouper(source Source, mode GroupMode, sep string, classifier func(string) string) (*Grouper, error) {
	if source == nil {
		return nil, errors.New("group: source must not be nil")
	}
	if mode == GroupByClassifier && classifier == nil {
		return nil, errors.New("group: classifier func required for GroupByClassifier mode")
	}
	if (mode == GroupByPrefix || mode == GroupBySuffix) && sep == "" {
		return nil, errors.New("group: separator must not be empty for prefix/suffix mode")
	}
	return &Grouper{source: source, mode: mode, sep: sep, classifier: classifier}, nil
}

// Group returns a map of group-name → map[key]value.
func (g *Grouper) Group() map[string]map[string]string {
	result := make(map[string]map[string]string)
	for _, key := range g.source.Keys() {
		val, _ := g.source.Get(key)
		groupName := g.classify(key)
		if result[groupName] == nil {
			result[groupName] = make(map[string]string)
		}
		result[groupName][key] = val
	}
	return result
}

// Keys returns all keys belonging to the named group.
func (g *Grouper) Keys(groupName string) []string {
	var out []string
	for _, key := range g.source.Keys() {
		if g.classify(key) == groupName {
			out = append(out, key)
		}
	}
	return out
}

func (g *Grouper) classify(key string) string {
	switch g.mode {
	case GroupByPrefix:
		if idx := strings.Index(key, g.sep); idx >= 0 {
			return key[:idx]
		}
		return ""
	case GroupBySuffix:
		if idx := strings.LastIndex(key, g.sep); idx >= 0 {
			return key[idx+len(g.sep):]
		}
		return ""
	case GroupByClassifier:
		return g.classifier(key)
	}
	return ""
}
