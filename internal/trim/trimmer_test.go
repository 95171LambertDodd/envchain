package trim_test

import (
	"testing"

	"github.com/yourorg/envchain/internal/trim"
)

// stubSource is a minimal in-memory Source for testing.
type stubSource struct {
	keys []string
	data map[string]string
}

func newStub(pairs ...string) *stubSource {
	s := &stubSource{data: make(map[string]string)}
	for i := 0; i+1 < len(pairs); i += 2 {
		s.keys = append(s.keys, pairs[i])
		s.data[pairs[i]] = pairs[i+1]
	}
	return s
}

func (s *stubSource) Keys() []string          { return s.keys }
func (s *stubSource) Get(k string) (string, bool) { v, ok := s.data[k]; return v, ok }

func TestNewTrimmerNilSourceError(t *testing.T) {
	_, err := trim.NewTrimmer(nil, trim.TrimBoth)
	if err == nil {
		t.Fatal("expected error for nil source")
	}
}

func TestTrimBothTrimsValues(t *testing.T) {
	src := newStub("KEY", "  hello  ", "OTHER", "\tworld\n")
	tr, _ := trim.NewTrimmer(src, trim.TrimBoth)
	out := tr.Apply()
	if out["KEY"] != "hello" {
		t.Errorf("expected 'hello', got %q", out["KEY"])
	}
	if out["OTHER"] != "world" {
		t.Errorf("expected 'world', got %q", out["OTHER"])
	}
}

func TestTrimKeysTrimsKeyNames(t *testing.T) {
	src := newStub("  SPACED  ", "value")
	tr, _ := trim.NewTrimmer(src, trim.TrimKeys)
	out := tr.Apply()
	if _, ok := out["SPACED"]; !ok {
		t.Error("expected trimmed key 'SPACED' to be present")
	}
	// value should be untouched
	if out["SPACED"] != "value" {
		t.Errorf("expected value 'value', got %q", out["SPACED"])
	}
}

func TestTrimAllTrimsBoth(t *testing.T) {
	src := newStub("  KEY  ", "  val  ")
	tr, _ := trim.NewTrimmer(src, trim.TrimAll)
	out := tr.Apply()
	if out["KEY"] != "val" {
		t.Errorf("expected 'val', got %q", out["KEY"])
	}
}

func TestTrimSkipsEmptyKeyAfterTrim(t *testing.T) {
	src := newStub("   ", "orphan")
	tr, _ := trim.NewTrimmer(src, trim.TrimKeys)
	out := tr.Apply()
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}

func TestApplyToLayer(t *testing.T) {
	src := newStub("A", "  trimmed  ", "B", " also ")
	dest := &collectLayer{data: make(map[string]string)}
	err := trim.ApplyToLayer(src, dest)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dest.data["A"] != "trimmed" {
		t.Errorf("expected 'trimmed', got %q", dest.data["A"])
	}
	if dest.data["B"] != "also" {
		t.Errorf("expected 'also', got %q", dest.data["B"])
	}
}

type collectLayer struct{ data map[string]string }

func (c *collectLayer) Set(k, v string) error { c.data[k] = v; return nil }
