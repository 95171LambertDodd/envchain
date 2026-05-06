package tag

import (
	"errors"
	"sort"
)

// FilterByTag returns a sorted slice of keys from the provided map whose
// entries carry ALL of the specified tags according to the given Tagger.
// Returns an error if tagger is nil or no tags are provided.
func FilterByTag(tagger *Tagger, entries map[string]string, tags ...string) ([]string, error) {
	if tagger == nil {
		return nil, errors.New("tag: tagger must not be nil")
	}
	if len(tags) == 0 {
		return nil, errors.New("tag: at least one tag is required")
	}

	var matched []string
	for key := range entries {
		if hasAllTags(tagger, key, tags) {
			matched = append(matched, key)
		}
	}
	sort.Strings(matched)
	return matched, nil
}

// hasAllTags reports whether the key carries every tag in the required list.
func hasAllTags(tagger *Tagger, key string, required []string) bool {
	keyTagSet := make(map[string]struct{})
	for _, t := range tagger.Tags(key) {
		keyTagSet[t] = struct{}{}
	}
	for _, req := range required {
		if _, ok := keyTagSet[req]; !ok {
			return false
		}
	}
	return true
}

// GroupByTag partitions the keys of entries into a map of tag -> []key.
// Keys with no tags are omitted. Each key may appear under multiple tags.
func GroupByTag(tagger *Tagger, entries map[string]string) map[string][]string {
	result := make(map[string][]string)
	for key := range entries {
		for _, tag := range tagger.Tags(key) {
			result[tag] = append(result[tag], key)
		}
	}
	for tag := range result {
		sort.Strings(result[tag])
	}
	return result
}
