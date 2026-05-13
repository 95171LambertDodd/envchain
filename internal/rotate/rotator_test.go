package rotate_test

import (
	"fmt"
	"testing"

	"github.com/your-org/envchain/internal/rotate"
)

// stubSource is an in-memory Source used in tests.
type stubSource struct {
	data map[string]string
}

func newStub(pairs ...string) *stubSource {
	data := make(map[string]string)
	for i := 0; i+1 < len(pairs); i += 2 {
		data[pairs[i]] = pairs[i+1]
	}
	return &stubSource{data: data}
}

func (s *stubSource) Get(key string) (string, bool) {
	v, ok := s.data[key]
	return v, ok
}

func (s *stubSource) Set(key, value string) error {
	if key == "" {
		return fmt.Errorf("empty key")
	}
	s.data[key] = value
	return nil
}

func (s *stubSource) Keys() []string {
	keys := make([]string, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}
	return keys
}

func TestNewRotatorNilSourceError(t *testing.T) {
	_, err := rotate.NewRotator(nil, nil)
	if err == nil {
		t.Fatal("expected error for nil source")
	}
}

func TestRotateSuccess(t *testing.T) {
	src := newStub("DB_PASS", "old")
	r, _ := rotate.NewRotator(src, nil)
	rec, err := r.Rotate("DB_PASS", "new")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.OldValue != "old" || rec.NewValue != "new" || rec.Key != "DB_PASS" {
		t.Errorf("unexpected record: %+v", rec)
	}
	if v, _ := src.Get("DB_PASS"); v != "new" {
		t.Errorf("expected source to be updated, got %q", v)
	}
}

func TestRotateMissingKeyError(t *testing.T) {
	src := newStub()
	r, _ := rotate.NewRotator(src, nil)
	_, err := r.Rotate("MISSING", "v")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestRotateNotAllowedError(t *testing.T) {
	src := newStub("SECRET", "val")
	r, _ := rotate.NewRotator(src, []string{"OTHER"})
	_, err := r.Rotate("SECRET", "new")
	if err == nil {
		t.Fatal("expected error when key not in allowed set")
	}
}

func TestRotateAllPartialErrors(t *testing.T) {
	src := newStub("A", "1", "B", "2")
	r, _ := rotate.NewRotator(src, nil)
	recs, errs := r.RotateAll(map[string]string{"A": "10", "MISSING": "99"})
	if len(recs) != 1 {
		t.Errorf("expected 1 record, got %d", len(recs))
	}
	if len(errs) != 1 {
		t.Errorf("expected 1 error, got %d", len(errs))
	}
}

func TestNewRotatorEmptyAllowedKeyError(t *testing.T) {
	_, err := rotate.NewRotator(newStub(), []string{""})
	if err == nil {
		t.Fatal("expected error for empty allowed key")
	}
}
