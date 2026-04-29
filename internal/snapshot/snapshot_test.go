package snapshot

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewSnapshotCopiesEntries(t *testing.T) {
	original := map[string]string{"KEY": "value"}
	s := NewSnapshot("test", original)
	original["KEY"] = "mutated"
	if s.Entries["KEY"] != "value" {
		t.Errorf("expected snapshot to be isolated from source map")
	}
	if s.Label != "test" {
		t.Errorf("expected label 'test', got %q", s.Label)
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	s := NewSnapshot("prod", map[string]string{"DB_URL": "postgres://localhost", "PORT": "5432"})
	if err := s.Save(path); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if loaded.Label != "prod" {
		t.Errorf("expected label 'prod', got %q", loaded.Label)
	}
	if loaded.Entries["PORT"] != "5432" {
		t.Errorf("expected PORT=5432, got %q", loaded.Entries["PORT"])
	}
}

func TestLoadInvalidPath(t *testing.T) {
	_, err := Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error loading nonexistent file")
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	os.WriteFile(path, []byte("not json"), 0644)
	_, err := Load(path)
	if err == nil {
		t.Error("expected error on invalid JSON")
	}
}

func TestDiff(t *testing.T) {
	base := NewSnapshot("base", map[string]string{"A": "1", "B": "2", "C": "3"})
	next := NewSnapshot("next", map[string]string{"A": "1", "B": "changed", "D": "4"})

	added, removed, changed := Diff(base, next)

	if len(added) != 1 || added[0] != "D" {
		t.Errorf("expected added=[D], got %v", added)
	}
	if len(removed) != 1 || removed[0] != "C" {
		t.Errorf("expected removed=[C], got %v", removed)
	}
	if len(changed) != 1 || changed[0] != "B" {
		t.Errorf("expected changed=[B], got %v", changed)
	}
}

func TestDiffNoChanges(t *testing.T) {
	entries := map[string]string{"X": "1"}
	base := NewSnapshot("a", entries)
	next := NewSnapshot("b", entries)
	added, removed, changed := Diff(base, next)
	if len(added)+len(removed)+len(changed) != 0 {
		t.Errorf("expected no diff, got added=%v removed=%v changed=%v", added, removed, changed)
	}
}
