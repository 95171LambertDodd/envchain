package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot captures the resolved state of a chain at a point in time.
type Snapshot struct {
	Timestamp time.Time         `json:"timestamp"`
	Label     string            `json:"label"`
	Entries   map[string]string `json:"entries"`
}

// NewSnapshot creates a new Snapshot with the given label and resolved entries.
func NewSnapshot(label string, entries map[string]string) *Snapshot {
	copy := make(map[string]string, len(entries))
	for k, v := range entries {
		copy[k] = v
	}
	return &Snapshot{
		Timestamp: time.Now().UTC(),
		Label:     label,
		Entries:   copy,
	}
}

// Save writes the snapshot as JSON to the given file path.
func (s *Snapshot) Save(path string) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal failed: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("snapshot: write failed: %w", err)
	}
	return nil
}

// Load reads a snapshot from the given JSON file path.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: read failed: %w", err)
	}
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal failed: %w", err)
	}
	return &s, nil
}

// Diff compares two snapshots and returns added, removed, and changed keys.
func Diff(base, next *Snapshot) (added, removed, changed []string) {
	for k, v := range next.Entries {
		if bv, ok := base.Entries[k]; !ok {
			added = append(added, k)
		} else if bv != v {
			changed = append(changed, k)
		}
	}
	for k := range base.Entries {
		if _, ok := next.Entries[k]; !ok {
			removed = append(removed, k)
		}
	}
	return
}
