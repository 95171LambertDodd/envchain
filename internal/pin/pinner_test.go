package pin_test

import (
	"errors"
	"testing"

	"github.com/yourorg/envchain/internal/pin"
)

type mapSource map[string]string

func (m mapSource) Get(key string) (string, bool) {
	v, ok := m[key]
	return v, ok
}

func TestPinEmptyKeyError(t *testing.T) {
	p := pin.NewPinner()
	if err := p.Pin("", "v"); err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestPinAndValidatePass(t *testing.T) {
	p := pin.NewPinner()
	_ = p.Pin("APP_ENV", "production")
	src := mapSource{"APP_ENV": "production"}
	if errs := p.Validate(src); len(errs) != 0 {
		t.Fatalf("expected no violations, got %v", errs)
	}
}

func TestPinViolationWrongValue(t *testing.T) {
	p := pin.NewPinner()
	_ = p.Pin("APP_ENV", "production")
	src := mapSource{"APP_ENV": "staging"}
	errs := p.Validate(src)
	if len(errs) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(errs))
	}
	if !errors.Is(errs[0], pin.ErrPinViolation) {
		t.Fatalf("expected ErrPinViolation, got %v", errs[0])
	}
}

func TestPinViolationMissingKey(t *testing.T) {
	p := pin.NewPinner()
	_ = p.Pin("DB_HOST", "localhost")
	src := mapSource{}
	errs := p.Validate(src)
	if len(errs) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(errs))
	}
	if !errors.Is(errs[0], pin.ErrPinViolation) {
		t.Fatalf("expected ErrPinViolation, got %v", errs[0])
	}
}

func TestUnpinRemovesKey(t *testing.T) {
	p := pin.NewPinner()
	_ = p.Pin("APP_ENV", "production")
	p.Unpin("APP_ENV")
	src := mapSource{"APP_ENV": "staging"}
	if errs := p.Validate(src); len(errs) != 0 {
		t.Fatalf("expected no violations after unpin, got %v", errs)
	}
}

func TestPinnedReturnsKeys(t *testing.T) {
	p := pin.NewPinner()
	_ = p.Pin("A", "1")
	_ = p.Pin("B", "2")
	keys := p.Pinned()
	if len(keys) != 2 {
		t.Fatalf("expected 2 pinned keys, got %d", len(keys))
	}
}

func TestMultipleViolations(t *testing.T) {
	p := pin.NewPinner()
	_ = p.Pin("A", "correct")
	_ = p.Pin("B", "correct")
	src := mapSource{"A": "wrong", "B": "alsowrong"}
	errs := p.Validate(src)
	if len(errs) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(errs))
	}
}
