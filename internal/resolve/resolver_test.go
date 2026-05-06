package resolve_test

import (
	"errors"
	"testing"

	"github.com/yourorg/envchain/internal/resolve"
)

// mapSource is a simple in-memory Source for testing.
type mapSource map[string]string

func (m mapSource) Get(key string) (string, bool) {
	v, ok := m[key]
	return v, ok
}

func makeResolver(layers ...struct {
	Name string
	Src  resolve.Source
}) *resolve.Resolver {
	return resolve.NewResolver(layers...)
}

func TestResolveFirstSourceWins(t *testing.T) {
	r := resolve.NewResolver(
		struct {
			Name string
			Src  resolve.Source
		}{"prod", mapSource{"DB_URL": "prod-db"}},
		struct {
			Name string
			Src  resolve.Source
		}{"dev", mapSource{"DB_URL": "dev-db"}},
	)
	res, err := r.Resolve("DB_URL")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Value != "prod-db" {
		t.Errorf("expected prod-db, got %q", res.Value)
	}
	if res.Source != "prod" {
		t.Errorf("expected source prod, got %q", res.Source)
	}
}

func TestResolveFallsThrough(t *testing.T) {
	r := resolve.NewResolver(
		struct {
			Name string
			Src  resolve.Source
		}{"prod", mapSource{}},
		struct {
			Name string
			Src  resolve.Source
		}{"dev", mapSource{"APP_ENV": "development"}},
	)
	res, err := r.Resolve("APP_ENV")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Value != "development" {
		t.Errorf("expected development, got %q", res.Value)
	}
	if res.Source != "dev" {
		t.Errorf("expected source dev, got %q", res.Source)
	}
}

func TestResolveMissingKey(t *testing.T) {
	r := resolve.NewResolver(
		struct {
			Name string
			Src  resolve.Source
		}{"base", mapSource{}},
	)
	_, err := r.Resolve("MISSING")
	if !errors.Is(err, resolve.ErrKeyNotFound) {
		t.Fatalf("expected ErrKeyNotFound, got %v", err)
	}
}

func TestResolveAllSuccess(t *testing.T) {
	r := resolve.NewResolver(
		struct {
			Name string
			Src  resolve.Source
		}{"base", mapSource{"A": "1", "B": "2"}},
	)
	results, err := r.ResolveAll([]string{"A", "B"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

func TestResolveAllPartialMissing(t *testing.T) {
	r := resolve.NewResolver(
		struct {
			Name string
			Src  resolve.Source
		}{"base", mapSource{"A": "1"}},
	)
	results, err := r.ResolveAll([]string{"A", "MISSING"})
	if !errors.Is(err, resolve.ErrKeyNotFound) {
		t.Fatalf("expected ErrKeyNotFound, got %v", err)
	}
	// partial results still returned for found keys
	if len(results) != 1 {
		t.Errorf("expected 1 partial result, got %d", len(results))
	}
}
