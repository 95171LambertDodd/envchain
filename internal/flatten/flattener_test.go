package flatten_test

import (
	"errors"
	"testing"

	"github.com/yourorg/envchain/internal/flatten"
)

// stubSource is a minimal in-memory Source for testing.
type stubSource struct {
	name    string
	entries map[string]string
	order   []string
}

func newStub(name string, kv map[string]string) *stubSource {
	order := make([]string, 0, len(kv))
	for k := range kv {
		order = append(order, k)
	}
	return &stubSource{name: name, entries: kv, order: order}
}

func (s *stubSource) Name() string           { return s.name }
func (s *stubSource) Keys() []string         { return s.order }
func (s *stubSource) Get(k string) (string, bool) {
	v, ok := s.entries[k]
	return v, ok
}

func TestFlattenLastWins(t *testing.T) {
	base := newStub("base", map[string]string{"A": "1", "B": "2"})
	over := newStub("over", map[string]string{"B": "99", "C": "3"})

	f := flatten.NewFlattener(flatten.StrategyLastWins)
	out, err := f.Flatten([]flatten.Source{base, over})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "1" || out["B"] != "99" || out["C"] != "3" {
		t.Errorf("unexpected map: %v", out)
	}
}

func TestFlattenFirstWins(t *testing.T) {
	base := newStub("base", map[string]string{"A": "original", "B": "2"})
	over := newStub("over", map[string]string{"A": "override", "C": "3"})

	f := flatten.NewFlattener(flatten.StrategyFirstWins)
	out, err := f.Flatten([]flatten.Source{base, over})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "original" {
		t.Errorf("expected 'original', got %q", out["A"])
	}
	if out["C"] != "3" {
		t.Errorf("expected C=3, got %q", out["C"])
	}
}

func TestFlattenErrorStrategyConflict(t *testing.T) {
	base := newStub("base", map[string]string{"X": "1"})
	over := newStub("over", map[string]string{"X": "2"})

	f := flatten.NewFlattener(flatten.StrategyError)
	_, err := f.Flatten([]flatten.Source{base, over})
	if !errors.Is(err, flatten.ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
}

func TestFlattenErrorStrategyNoConflict(t *testing.T) {
	base := newStub("base", map[string]string{"A": "1"})
	over := newStub("over", map[string]string{"B": "2"})

	f := flatten.NewFlattener(flatten.StrategyError)
	out, err := f.Flatten([]flatten.Source{base, over})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
}

func TestFlattenWithPrefix(t *testing.T) {
	src := newStub("src", map[string]string{"HOST": "localhost", "PORT": "5432"})

	f := flatten.NewFlattener(flatten.StrategyLastWins).WithPrefix("DB_")
	out, err := f.Flatten([]flatten.Source{src})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_HOST"] != "localhost" || out["DB_PORT"] != "5432" {
		t.Errorf("unexpected prefixed map: %v", out)
	}
	if _, bare := out["HOST"]; bare {
		t.Error("expected no unprefixed key HOST")
	}
}
