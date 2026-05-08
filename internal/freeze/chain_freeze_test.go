package freeze_test

import (
	"errors"
	"testing"

	"github.com/yourorg/envchain/internal/freeze"
)

// stubSource implements freeze.Source for testing.
type stubSource struct {
	entries map[string]string
	order   []string
}

func newStub(pairs ...string) *stubSource {
	s := &stubSource{entries: make(map[string]string)}
	for i := 0; i+1 < len(pairs); i += 2 {
		s.order = append(s.order, pairs[i])
		s.entries[pairs[i]] = pairs[i+1]
	}
	return s
}

func (s *stubSource) Keys() []string        { return s.order }
func (s *stubSource) Get(k string) (string, bool) { v, ok := s.entries[k]; return v, ok }

func TestFreezeFromSourceLocksAllKeys(t *testing.T) {
	f := freeze.NewFreezer()
	src := newStub("HOST", "localhost", "PORT", "5432")
	if err := freeze.FreezeFromSource(f, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, k := range []string{"HOST", "PORT"} {
		if !f.IsFrozen(k) {
			t.Errorf("expected %s to be frozen", k)
		}
	}
}

func TestFreezeFromSourceValuesPreserved(t *testing.T) {
	f := freeze.NewFreezer()
	src := newStub("ENV", "production")
	_ = freeze.FreezeFromSource(f, src)
	v, ok := f.Get("ENV")
	if !ok || v != "production" {
		t.Fatalf("expected ENV=production, got %q ok=%v", v, ok)
	}
}

func TestFreezeFromSourcePreventsOverride(t *testing.T) {
	f := freeze.NewFreezer()
	src := newStub("SECRET", "abc123")
	_ = freeze.FreezeFromSource(f, src)
	err := f.Set("SECRET", "hacked")
	if !errors.Is(err, freeze.ErrFrozen) {
		t.Fatalf("expected ErrFrozen, got %v", err)
	}
}

func TestFreezeFromSourceEmptySource(t *testing.T) {
	f := freeze.NewFreezer()
	src := newStub()
	if err := freeze.FreezeFromSource(f, src); err != nil {
		t.Fatalf("unexpected error on empty source: %v", err)
	}
	if len(f.FrozenKeys()) != 0 {
		t.Fatal("expected no frozen keys for empty source")
	}
}
