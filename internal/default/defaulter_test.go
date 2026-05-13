package defaulter_test

import (
	"testing"

	defaulter "github.com/yourorg/envchain/internal/default"
)

// stubLayer is a minimal in-memory key/value store used as both Source
// and Target in tests.
type stubLayer struct {
	data map[string]string
}

func newStub(pairs ...string) *stubLayer {
	s := &stubLayer{data: make(map[string]string)}
	for i := 0; i+1 < len(pairs); i += 2 {
		s.data[pairs[i]] = pairs[i+1]
	}
	return s
}

func (s *stubLayer) Keys() []string {
	keys := make([]string, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}
	return keys
}

func (s *stubLayer) Get(key string) (string, bool) {
	v, ok := s.data[key]
	return v, ok
}

func (s *stubLayer) Set(key, value string) error {
	if key == "" {
		return fmt.Errorf("empty key")
	}
	s.data[key] = value
	return nil
}

func TestNewDefaulterEmptyMapError(t *testing.T) {
	_, err := defaulter.NewDefaulter(nil)
	if err == nil {
		t.Fatal("expected error for nil defaults")
	}
}

func TestNewDefaulterEmptyKeyError(t *testing.T) {
	_, err := defaulter.NewDefaulter(map[string]string{"": "value"})
	if err == nil {
		t.Fatal("expected error for empty key in defaults")
	}
}

func TestApplyFillsMissingKey(t *testing.T) {
	d, _ := defaulter.NewDefaulter(map[string]string{"HOST": "localhost", "PORT": "8080"})
	target := newStub()
	if err := d.Apply(target); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, _ := target.Get("HOST"); v != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", v)
	}
	if v, _ := target.Get("PORT"); v != "8080" {
		t.Errorf("expected PORT=8080, got %q", v)
	}
}

func TestApplyDoesNotOverwriteExistingValue(t *testing.T) {
	d, _ := defaulter.NewDefaulter(map[string]string{"HOST": "localhost"})
	target := newStub("HOST", "production.example.com")
	if err := d.Apply(target); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, _ := target.Get("HOST"); v != "production.example.com" {
		t.Errorf("expected original value preserved, got %q", v)
	}
}

func TestApplyOverwritesEmptyValue(t *testing.T) {
	d, _ := defaulter.NewDefaulter(map[string]string{"LOG_LEVEL": "info"})
	target := newStub("LOG_LEVEL", "")
	if err := d.Apply(target); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, _ := target.Get("LOG_LEVEL"); v != "info" {
		t.Errorf("expected LOG_LEVEL=info, got %q", v)
	}
}

func TestApplyFromSourceFillsMissing(t *testing.T) {
	d, _ := defaulter.NewDefaulter(map[string]string{"DUMMY": "x"})
	src := newStub("DB_HOST", "db.local", "DB_PORT", "5432")
	target := newStub("DB_HOST", "override.local")
	if err := d.ApplyFromSource(src, target); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, _ := target.Get("DB_HOST"); v != "override.local" {
		t.Errorf("expected DB_HOST unchanged, got %q", v)
	}
	if v, _ := target.Get("DB_PORT"); v != "5432" {
		t.Errorf("expected DB_PORT=5432 from source, got %q", v)
	}
}
