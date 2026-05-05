package scope_test

import (
	"errors"
	"testing"

	"github.com/yourorg/envchain/internal/scope"
)

func TestNewScopeEmptyNameError(t *testing.T) {
	_, err := scope.NewScope("")
	if err == nil {
		t.Fatal("expected error for empty scope name")
	}
}

func TestNewScopeSuccess(t *testing.T) {
	s, err := scope.NewScope("dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Name() != "dev" {
		t.Errorf("expected name %q, got %q", "dev", s.Name())
	}
}

func TestScopeSetAndGet(t *testing.T) {
	s, _ := scope.NewScope("staging")
	if err := s.Set("DB_HOST", "db.staging.local"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, ok := s.Get("DB_HOST")
	if !ok || v != "db.staging.local" {
		t.Errorf("expected %q, got %q (ok=%v)", "db.staging.local", v, ok)
	}
}

func TestScopeSetEmptyKeyError(t *testing.T) {
	s, _ := scope.NewScope("prod")
	if err := s.Set("", "value"); err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestScopeGetMissingKey(t *testing.T) {
	s, _ := scope.NewScope("dev")
	_, ok := s.Get("MISSING")
	if ok {
		t.Fatal("expected missing key to return false")
	}
}

func TestRegistryActivateUnknownScope(t *testing.T) {
	r := scope.NewScopeRegistry()
	err := r.Activate("ghost")
	if !errors.Is(err, scope.ErrUnknownScope) {
		t.Errorf("expected ErrUnknownScope, got %v", err)
	}
}

func TestRegistryActivateAndResolve(t *testing.T) {
	r := scope.NewScopeRegistry()

	dev, _ := scope.NewScope("dev")
	_ = dev.Set("API_URL", "http://localhost")

	prod, _ := scope.NewScope("prod")
	_ = prod.Set("API_URL", "https://api.example.com")

	r.Register(dev)
	r.Register(prod)

	_ = r.Activate("dev")
	v, ok := r.Resolve("API_URL")
	if !ok || v != "http://localhost" {
		t.Errorf("dev: expected %q, got %q", "http://localhost", v)
	}

	_ = r.Activate("prod")
	v, ok = r.Resolve("API_URL")
	if !ok || v != "https://api.example.com" {
		t.Errorf("prod: expected %q, got %q", "https://api.example.com", v)
	}
}

func TestRegistryResolveNoActiveScope(t *testing.T) {
	r := scope.NewScopeRegistry()
	_, ok := r.Resolve("ANY_KEY")
	if ok {
		t.Fatal("expected false when no scope is active")
	}
}
