package group_test

import (
	"testing"

	"github.com/yourorg/envchain/internal/group"
)

// TestGroupEmptySource verifies that grouping an empty source yields an empty map.
func TestGroupEmptySource(t *testing.T) {
	src := newStub()
	g, err := group.NewGrouper(src, group.GroupByPrefix, "_", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := g.Group()
	if len(result) != 0 {
		t.Errorf("expected empty groups, got %d", len(result))
	}
}

// TestGroupNoSeparatorMatchFallsToEmpty ensures keys with no separator land in the "" group.
func TestGroupNoSeparatorMatchFallsToEmpty(t *testing.T) {
	src := newStub("NOSEP", "value")
	g, err := group.NewGrouper(src, group.GroupByPrefix, "_", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := g.Group()
	if _, ok := result[""]; !ok {
		t.Error("expected ungrouped keys under empty-string group")
	}
	if result[""]["NOSEP"] != "value" {
		t.Errorf("expected value 'value', got %q", result[""]["NOSEP"])
	}
}

// TestGroupValuesPreserved checks that the original values are intact after grouping.
func TestGroupValuesPreserved(t *testing.T) {
	src := newStub(
		"DB_HOST", "db.internal",
		"DB_PASS", "s3cr3t",
	)
	g, err := group.NewGrouper(src, group.GroupByPrefix, "_", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	groups := g.Group()
	if groups["DB"]["DB_HOST"] != "db.internal" {
		t.Errorf("unexpected value for DB_HOST: %q", groups["DB"]["DB_HOST"])
	}
	if groups["DB"]["DB_PASS"] != "s3cr3t" {
		t.Errorf("unexpected value for DB_PASS: %q", groups["DB"]["DB_PASS"])
	}
}

// TestGroupMutationDoesNotAffectSource ensures the returned map is a copy.
func TestGroupMutationDoesNotAffectSource(t *testing.T) {
	src := newStub("APP_NAME", "envchain")
	g, err := group.NewGrouper(src, group.GroupByPrefix, "_", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	groups := g.Group()
	groups["APP"]["APP_NAME"] = "mutated"

	// Re-group to confirm source is untouched.
	groups2 := g.Group()
	if groups2["APP"]["APP_NAME"] != "envchain" {
		t.Errorf("source was mutated; got %q", groups2["APP"]["APP_NAME"])
	}
}

// TestGroupKeysUnknownGroup returns nil for a group that does not exist.
func TestGroupKeysUnknownGroup(t *testing.T) {
	src := newStub("DB_HOST", "localhost")
	g, err := group.NewGrouper(src, group.GroupByPrefix, "_", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	keys := g.Keys("NONEXISTENT")
	if len(keys) != 0 {
		t.Errorf("expected no keys for unknown group, got %d", len(keys))
	}
}
