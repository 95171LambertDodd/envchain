package group_test

import (
	"sort"
	"strings"
	"testing"

	"github.com/yourorg/envchain/internal/group"
)

// stubSource is a minimal in-memory key-value store for testing.
type stubSource struct {
	data map[string]string
}

func newStub(pairs ...string) *stubSource {
	s := &stubSource{data: make(map[string]string)}
	for i := 0; i+1 < len(pairs); i += 2 {
		s.data[pairs[i]] = pairs[i+1]
	}
	return s
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

func TestNewGrouperNilSourceError(t *testing.T) {
	_, err := group.NewGrouper(nil, group.GroupByPrefix, "_", nil)
	if err == nil {
		t.Fatal("expected error for nil source")
	}
}

func TestNewGrouperMissingSeparator(t *testing.T) {
	src := newStub()
	_, err := group.NewGrouper(src, group.GroupByPrefix, "", nil)
	if err == nil {
		t.Fatal("expected error for empty separator")
	}
}

func TestNewGrouperMissingClassifier(t *testing.T) {
	src := newStub()
	_, err := group.NewGrouper(src, group.GroupByClassifier, "", nil)
	if err == nil {
		t.Fatal("expected error for nil classifier")
	}
}

func TestGroupByPrefix(t *testing.T) {
	src := newStub(
		"DB_HOST", "localhost",
		"DB_PORT", "5432",
		"APP_NAME", "envchain",
	)
	g, err := group.NewGrouper(src, group.GroupByPrefix, "_", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	groups := g.Group()
	if len(groups["DB"]) != 2 {
		t.Errorf("expected 2 DB keys, got %d", len(groups["DB"]))
	}
	if len(groups["APP"]) != 1 {
		t.Errorf("expected 1 APP key, got %d", len(groups["APP"]))
	}
}

func TestGroupBySuffix(t *testing.T) {
	src := newStub(
		"host_dev", "localhost",
		"host_prod", "prod.example.com",
		"port_dev", "5432",
	)
	g, err := group.NewGrouper(src, group.GroupBySuffix, "_", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	groups := g.Group()
	if len(groups["dev"]) != 2 {
		t.Errorf("expected 2 dev keys, got %d", len(groups["dev"]))
	}
	if len(groups["prod"]) != 1 {
		t.Errorf("expected 1 prod key, got %d", len(groups["prod"]))
	}
}

func TestGroupByClassifier(t *testing.T) {
	src := newStub(
		"SECRET_KEY", "abc",
		"PUBLIC_URL", "https://example.com",
		"SECRET_TOKEN", "xyz",
	)
	classifier := func(key string) string {
		if strings.HasPrefix(key, "SECRET") {
			return "sensitive"
		}
		return "public"
	}
	g, err := group.NewGrouper(src, group.GroupByClassifier, "", classifier)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	groups := g.Group()
	if len(groups["sensitive"]) != 2 {
		t.Errorf("expected 2 sensitive keys, got %d", len(groups["sensitive"]))
	}
	if len(groups["public"]) != 1 {
		t.Errorf("expected 1 public key, got %d", len(groups["public"]))
	}
}

func TestGroupKeysMethod(t *testing.T) {
	src := newStub(
		"DB_HOST", "localhost",
		"DB_PORT", "5432",
		"APP_NAME", "envchain",
	)
	g, _ := group.NewGrouper(src, group.GroupByPrefix, "_", nil)
	keys := g.Keys("DB")
	if len(keys) != 2 {
		t.Errorf("expected 2 keys for group DB, got %d", len(keys))
	}
}
