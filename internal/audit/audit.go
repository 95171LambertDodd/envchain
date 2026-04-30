// Package audit provides change tracking for environment config layers.
package audit

import (
	"fmt"
	"time"
)

// ChangeKind describes the type of change recorded.
type ChangeKind string

const (
	ChangeSet    ChangeKind = "set"
	ChangeDelete ChangeKind = "delete"
	ChangeMerge  ChangeKind = "merge"
)

// Entry represents a single audited change.
type Entry struct {
	Timestamp time.Time  `json:"timestamp"`
	Layer     string     `json:"layer"`
	Key       string     `json:"key"`
	OldValue  string     `json:"old_value,omitempty"`
	NewValue  string     `json:"new_value,omitempty"`
	Kind      ChangeKind `json:"kind"`
}

// Log holds an ordered list of audit entries.
type Log struct {
	entries []Entry
	clock   func() time.Time
}

// NewLog creates a new audit Log. Optionally inject a clock for testing.
func NewLog(clock func() time.Time) *Log {
	if clock == nil {
		clock = time.Now
	}
	return &Log{clock: clock}
}

// Record appends a new entry to the log.
func (l *Log) Record(layer, key, oldVal, newVal string, kind ChangeKind) {
	l.entries = append(l.entries, Entry{
		Timestamp: l.clock(),
		Layer:     layer,
		Key:       key,
		OldValue:  oldVal,
		NewValue:  newVal,
		Kind:      kind,
	})
}

// Entries returns a copy of all recorded entries.
func (l *Log) Entries() []Entry {
	out := make([]Entry, len(l.entries))
	copy(out, l.entries)
	return out
}

// FilterByLayer returns entries that belong to the given layer name.
func (l *Log) FilterByLayer(layer string) []Entry {
	var out []Entry
	for _, e := range l.entries {
		if e.Layer == layer {
			out = append(out, e)
		}
	}
	return out
}

// FilterByKey returns entries that match the given key.
func (l *Log) FilterByKey(key string) []Entry {
	var out []Entry
	for _, e := range l.entries {
		if e.Key == key {
			out = append(out, e)
		}
	}
	return out
}

// Summary returns a human-readable summary of all entries.
func (l *Log) Summary() string {
	if len(l.entries) == 0 {
		return "no audit entries"
	}
	var s string
	for _, e := range l.entries {
		s += fmt.Sprintf("[%s] %s layer=%s key=%s\n",
			e.Timestamp.Format(time.RFC3339), e.Kind, e.Layer, e.Key)
	}
	return s
}
