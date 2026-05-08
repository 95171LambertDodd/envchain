package alias_test

import (
	"testing"

	"github.com/yourorg/envchain/internal/alias"
)

// mapSource is a simple Source backed by a map.
type mapSource map[string]string

func (m mapSource) Get(key string) (string, bool) {
	v, ok := m[key]
	return v, ok
}

func TestNewAliasNilSourceError(t *testing.T) {
	_, err := alias.NewAliaser(nil)
	if err == nil {
		t.Fatal("expected error for nil source")
	}
}

func TestAddEmptyAliasError(t *testing.T) {
	a, _ := alias.NewAliaser(mapSource{})
	if err := a.Add("", "CANONICAL"); err == nil {
		t.Fatal("expected error for empty alias")
	}
}

func TestAddEmptyCanonicalError(t *testing.T) {
	a, _ := alias.NewAliaser(mapSource{})
	if err := a.Add("ALIAS", ""); err == nil {
		t.Fatal("expected error for empty canonical")
	}
}

func TestAddConflictError(t *testing.T) {
	a, _ := alias.NewAliaser(mapSource{})
	_ = a.Add("OLD_KEY", "NEW_KEY")
	if err := a.Add("OLD_KEY", "OTHER_KEY"); err == nil {
		t.Fatal("expected error for conflicting alias")
	}
}

func TestAddIdempotent(t *testing.T) {
	a, _ := alias.NewAliaser(mapSource{})
	_ = a.Add("OLD_KEY", "NEW_KEY")
	if err := a.Add("OLD_KEY", "NEW_KEY"); err != nil {
		t.Fatalf("unexpected error for idempotent add: %v", err)
	}
}

func TestResolveDirectKey(t *testing.T) {
	src := mapSource{"DATABASE_URL": "postgres://localhost/db"}
	a, _ := alias.NewAliaser(src)
	v, ok := a.Resolve("DATABASE_URL")
	if !ok || v != "postgres://localhost/db" {
		t.Fatalf("expected direct key resolution, got %q %v", v, ok)
	}
}

func TestResolveViaAlias(t *testing.T) {
	src := mapSource{"DATABASE_URL": "postgres://localhost/db"}
	a, _ := alias.NewAliaser(src)
	_ = a.Add("DB_URL", "DATABASE_URL")
	v, ok := a.Resolve("DB_URL")
	if !ok || v != "postgres://localhost/db" {
		t.Fatalf("expected alias resolution, got %q %v", v, ok)
	}
}

func TestResolveMissingKey(t *testing.T) {
	a, _ := alias.NewAliaser(mapSource{})
	_, ok := a.Resolve("MISSING")
	if ok {
		t.Fatal("expected false for missing key")
	}
}

func TestCanonicalWithAlias(t *testing.T) {
	a, _ := alias.NewAliaser(mapSource{})
	_ = a.Add("DB_URL", "DATABASE_URL")
	if got := a.Canonical("DB_URL"); got != "DATABASE_URL" {
		t.Fatalf("expected DATABASE_URL, got %q", got)
	}
}

func TestCanonicalWithoutAlias(t *testing.T) {
	a, _ := alias.NewAliaser(mapSource{})
	if got := a.Canonical("SOME_KEY"); got != "SOME_KEY" {
		t.Fatalf("expected SOME_KEY, got %q", got)
	}
}
