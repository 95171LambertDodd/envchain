package index_test

import (
	"testing"

	"github.com/yourorg/envchain/internal/index"
)

// stubSource is a minimal in-memory Source for testing.
type stubSource struct {
	data map[string]string
}

func (s *stubSource) Keys() []string {
	keys := make([]string, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}
	return keys
}

func (s *stubSource) Get(key string) (string, bool) {
	v, ok := s.data[key]
	return v, ok
}

func newStub(pairs ...string) *stubSource {
	m := make(map[string]string)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return &stubSource{data: m}
}

func TestNewIndexerNilSourceError(t *testing.T) {
	_, err := index.NewIndexer(nil)
	if err == nil {
		t.Fatal("expected error for nil source")
	}
}

func TestIndexerGetFound(t *testing.T) {
	idx, err := index.NewIndexer(newStub("FOO", "bar"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, ok := idx.Get("FOO")
	if !ok || v != "bar" {
		t.Fatalf("expected bar, got %q ok=%v", v, ok)
	}
}

func TestIndexerGetMissing(t *testing.T) {
	idx, _ := index.NewIndexer(newStub("A", "1"))
	_, ok := idx.Get("MISSING")
	if ok {
		t.Fatal("expected miss for unknown key")
	}
}

func TestIndexerKeysAreSorted(t *testing.T) {
	idx, _ := index.NewIndexer(newStub("Z", "z", "A", "a", "M", "m"))
	keys := idx.Keys()
	if len(keys) != 3 || keys[0] != "A" || keys[1] != "M" || keys[2] != "Z" {
		t.Fatalf("unexpected key order: %v", keys)
	}
}

func TestKeysForValueSingleMatch(t *testing.T) {
	idx, _ := index.NewIndexer(newStub("HOST", "localhost", "PORT", "5432"))
	keys := idx.KeysForValue("localhost")
	if len(keys) != 1 || keys[0] != "HOST" {
		t.Fatalf("unexpected keys: %v", keys)
	}
}

func TestKeysForValueMultipleMatch(t *testing.T) {
	idx, _ := index.NewIndexer(newStub("A", "same", "B", "same", "C", "other"))
	keys := idx.KeysForValue("same")
	if len(keys) != 2 || keys[0] != "A" || keys[1] != "B" {
		t.Fatalf("unexpected keys: %v", keys)
	}
}

func TestKeysForValueMissing(t *testing.T) {
	idx, _ := index.NewIndexer(newStub("X", "y"))
	keys := idx.KeysForValue("nope")
	if len(keys) != 0 {
		t.Fatalf("expected empty slice, got %v", keys)
	}
}

func TestHasDuplicateValuesTrue(t *testing.T) {
	idx, _ := index.NewIndexer(newStub("P", "dup", "Q", "dup"))
	ok, msg := idx.HasDuplicateValues()
	if !ok {
		t.Fatal("expected duplicates to be detected")
	}
	if msg == "" {
		t.Fatal("expected non-empty diagnostic message")
	}
}

func TestHasDuplicateValuesFalse(t *testing.T) {
	idx, _ := index.NewIndexer(newStub("A", "1", "B", "2", "C", "3"))
	ok, _ := idx.HasDuplicateValues()
	if ok {
		t.Fatal("expected no duplicates")
	}
}
