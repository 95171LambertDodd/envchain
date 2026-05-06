package tag

import (
	"errors"
	"sort"
	"strings"
)

// Tagger associates string tags with environment keys, enabling grouping,
// filtering, and annotation of config entries.
type Tagger struct {
	// keyTags maps a key to its set of tags.
	keyTags map[string]map[string]struct{}
}

// NewTagger returns an initialised Tagger.
func NewTagger() *Tagger {
	return &Tagger{
		keyTags: make(map[string]map[string]struct{}),
	}
}

// Tag adds one or more tags to the given key.
// Returns an error if key or any tag is empty.
func (t *Tagger) Tag(key string, tags ...string) error {
	if strings.TrimSpace(key) == "" {
		return errors.New("tag: key must not be empty")
	}
	for _, tag := range tags {
		if strings.TrimSpace(tag) == "" {
			return errors.New("tag: tag value must not be empty")
		}
	}
	if _, ok := t.keyTags[key]; !ok {
		t.keyTags[key] = make(map[string]struct{})
	}
	for _, tag := range tags {
		t.keyTags[key][tag] = struct{}{}
	}
	return nil
}

// Tags returns a sorted slice of tags associated with key.
// Returns nil if the key has no tags.
func (t *Tagger) Tags(key string) []string {
	set, ok := t.keyTags[key]
	if !ok || len(set) == 0 {
		return nil
	}
	out := make([]string, 0, len(set))
	for tag := range set {
		out = append(out, tag)
	}
	sort.Strings(out)
	return out
}

// KeysWithTag returns a sorted slice of all keys that carry the given tag.
func (t *Tagger) KeysWithTag(tag string) []string {
	var out []string
	for key, set := range t.keyTags {
		if _, ok := set[tag]; ok {
			out = append(out, key)
		}
	}
	sort.Strings(out)
	return out
}

// Remove removes a tag from a key. No-op if key or tag is absent.
func (t *Tagger) Remove(key, tag string) {
	if set, ok := t.keyTags[key]; ok {
		delete(set, tag)
	}
}

// AllTags returns a sorted deduplicated slice of every tag in use.
func (t *Tagger) AllTags() []string {
	set := make(map[string]struct{})
	for _, tags := range t.keyTags {
		for tag := range tags {
			set[tag] = struct{}{}
		}
	}
	out := make([]string, 0, len(set))
	for tag := range set {
		out = append(out, tag)
	}
	sort.Strings(out)
	return out
}
