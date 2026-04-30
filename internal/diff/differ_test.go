package diff_test

import (
	"testing"

	"github.com/your-org/envchain/internal/diff"
)

// mapResolver is a simple Resolver backed by a map.
type mapResolver struct {
	data map[string]string
}

func newMap(data map[string]string) *mapResolver {
	return &mapResolver{data: data}
}

func (m *mapResolver) Get(key string) (string, bool) {
	v, ok := m.data[key]
	return v, ok
}

func (m *mapResolver) Keys() []string {
	keys := make([]string, 0, len(m.data))
	for k := range m.data {
		keys = append(keys, k)
	}
	return keys
}

func TestDiffNoChanges(t *testing.T) {
	d := diff.NewDiffer()
	base := newMap(map[string]string{"A": "1", "B": "2"})
	changes := d.Diff(base, base)
	if len(changes) != 0 {
		t.Fatalf("expected no changes, got %d", len(changes))
	}
}

func TestDiffAdded(t *testing.T) {
	d := diff.NewDiffer()
	base := newMap(map[string]string{})
	head := newMap(map[string]string{"NEW_KEY": "hello"})
	changes := d.Diff(base, head)
	if len(changes) != 1 || changes[0].Kind != diff.Added || changes[0].Key != "NEW_KEY" {
		t.Fatalf("unexpected changes: %+v", changes)
	}
}

func TestDiffRemoved(t *testing.T) {
	d := diff.NewDiffer()
	base := newMap(map[string]string{"GONE": "bye"})
	head := newMap(map[string]string{})
	changes := d.Diff(base, head)
	if len(changes) != 1 || changes[0].Kind != diff.Removed || changes[0].OldValue != "bye" {
		t.Fatalf("unexpected changes: %+v", changes)
	}
}

func TestDiffModified(t *testing.T) {
	d := diff.NewDiffer()
	base := newMap(map[string]string{"PORT": "8080"})
	head := newMap(map[string]string{"PORT": "9090"})
	changes := d.Diff(base, head)
	if len(changes) != 1 || changes[0].Kind != diff.Modified {
		t.Fatalf("unexpected changes: %+v", changes)
	}
	if changes[0].OldValue != "8080" || changes[0].NewValue != "9090" {
		t.Fatalf("wrong values: %+v", changes[0])
	}
}

func TestDiffSorted(t *testing.T) {
	d := diff.NewDiffer()
	base := newMap(map[string]string{"Z": "1", "A": "1"})
	head := newMap(map[string]string{"Z": "2", "A": "2"})
	changes := d.Diff(base, head)
	if len(changes) != 2 || changes[0].Key != "A" || changes[1].Key != "Z" {
		t.Fatalf("changes not sorted: %+v", changes)
	}
}

func TestChangeString(t *testing.T) {
	cases := []struct {
		c    diff.Change
		want string
	}{
		{diff.Change{Key: "X", Kind: diff.Added, NewValue: "v"}, `+ X="v"`},
		{diff.Change{Key: "X", Kind: diff.Removed, OldValue: "v"}, `- X="v"`},
		{diff.Change{Key: "X", Kind: diff.Modified, OldValue: "a", NewValue: "b"}, `~ X: "a" -> "b"`},
	}
	for _, tc := range cases {
		if got := tc.c.String(); got != tc.want {
			t.Errorf("String() = %q, want %q", got, tc.want)
		}
	}
}
