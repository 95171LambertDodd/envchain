package namespace_test

import (
	"sort"
	"testing"

	"github.com/yourorg/envchain/internal/namespace"
)

// stubSource is a simple in-memory Source for testing.
type stubSource struct {
	data map[string]string
}

func (s *stubSource) Keys() []string {
	keys := make([]string, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (s *stubSource) Get(key string) (string, bool) {
	v, ok := s.data[key]
	return v, ok
}

func newStub(data map[string]string) *stubSource { return &stubSource{data: data} }

func TestNewNamespacerNilSourceError(t *testing.T) {
	_, err := namespace.NewNamespacer(nil, "prod", "_")
	if err == nil {
		t.Fatal("expected error for nil source")
	}
}

func TestNewNamespacerEmptyNamespaceError(t *testing.T) {
	_, err := namespace.NewNamespacer(newStub(nil), "", "_")
	if err == nil {
		t.Fatal("expected error for empty namespace")
	}
}

func TestNewNamespacerEmptySepError(t *testing.T) {
	_, err := namespace.NewNamespacer(newStub(nil), "prod", "")
	if err == nil {
		t.Fatal("expected error for empty separator")
	}
}

func TestQualifiedKey(t *testing.T) {
	n, _ := namespace.NewNamespacer(newStub(nil), "prod", "_")
	if got := n.QualifiedKey("DB_HOST"); got != "prod_DB_HOST" {
		t.Fatalf("expected prod_DB_HOST, got %s", got)
	}
}

func TestStrip(t *testing.T) {
	n, _ := namespace.NewNamespacer(newStub(nil), "prod", "_")

	bare, ok := n.Strip("prod_DB_HOST")
	if !ok || bare != "DB_HOST" {
		t.Fatalf("expected (DB_HOST, true), got (%s, %v)", bare, ok)
	}

	_, ok = n.Strip("staging_DB_HOST")
	if ok {
		t.Fatal("expected false for key from different namespace")
	}
}

func TestKeys(t *testing.T) {
	src := newStub(map[string]string{"FOO": "1", "BAR": "2"})
	n, _ := namespace.NewNamespacer(src, "app", ".")

	keys := n.Keys()
	sort.Strings(keys)
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
	if keys[0] != "app.BAR" || keys[1] != "app.FOO" {
		t.Fatalf("unexpected keys: %v", keys)
	}
}

func TestGet(t *testing.T) {
	src := newStub(map[string]string{"PORT": "8080"})
	n, _ := namespace.NewNamespacer(src, "svc", "__")

	v, ok := n.Get("svc__PORT")
	if !ok || v != "8080" {
		t.Fatalf("expected (8080, true), got (%s, %v)", v, ok)
	}

	_, ok = n.Get("other__PORT")
	if ok {
		t.Fatal("expected false for wrong namespace")
	}
}

func TestFlatten(t *testing.T) {
	src := newStub(map[string]string{"A": "1", "B": "2"})
	n, _ := namespace.NewNamespacer(src, "ns", "-")

	flat := n.Flatten()
	if len(flat) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(flat))
	}
	if flat["ns-A"] != "1" || flat["ns-B"] != "2" {
		t.Fatalf("unexpected flatten result: %v", flat)
	}
}
