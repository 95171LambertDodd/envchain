package config

import (
	"testing"
)

func TestLayerSetAndGet(t *testing.T) {
	l := NewLayer("base", EnvDev)
	if err := l.Set("APP_PORT", "8080"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, ok := l.Get("APP_PORT")
	if !ok || v != "8080" {
		t.Errorf("expected APP_PORT=8080, got %q ok=%v", v, ok)
	}
}

func TestLayerSetEmptyKeyError(t *testing.T) {
	l := NewLayer("base", EnvDev)
	if err := l.Set("", "value"); err == nil {
		t.Error("expected error for empty key, got nil")
	}
}

func TestChainResolveOverride(t *testing.T) {
	base := NewLayer("base", EnvDev)
	_ = base.Set("DB_HOST", "localhost")
	_ = base.Set("LOG_LEVEL", "debug")

	override := NewLayer("prod-override", EnvProd)
	_ = override.Set("DB_HOST", "db.prod.internal")

	chain := NewChain()
	chain.Push(base)
	chain.Push(override)

	resolved := chain.Resolve()
	if resolved["DB_HOST"] != "db.prod.internal" {
		t.Errorf("expected overridden DB_HOST, got %q", resolved["DB_HOST"])
	}
	if resolved["LOG_LEVEL"] != "debug" {
		t.Errorf("expected LOG_LEVEL=debug from base, got %q", resolved["LOG_LEVEL"])
	}
}

func TestChainGetTopWins(t *testing.T) {
	l1 := NewLayer("l1", EnvDev)
	_ = l1.Set("KEY", "first")
	l2 := NewLayer("l2", EnvStaging)
	_ = l2.Set("KEY", "second")

	chain := NewChain()
	chain.Push(l1)
	chain.Push(l2)

	v, ok := chain.Get("KEY")
	if !ok || v != "second" {
		t.Errorf("expected second, got %q ok=%v", v, ok)
	}
}

func TestChainRequireKeysMissing(t *testing.T) {
	l := NewLayer("base", EnvDev)
	_ = l.Set("PRESENT", "yes")

	chain := NewChain()
	chain.Push(l)

	err := chain.RequireKeys([]string{"PRESENT", "MISSING_ONE", "MISSING_TWO"})
	if err == nil {
		t.Fatal("expected error for missing keys, got nil")
	}
}

func TestChainRequireKeysAllPresent(t *testing.T) {
	l := NewLayer("base", EnvProd)
	_ = l.Set("A", "1")
	_ = l.Set("B", "2")

	chain := NewChain()
	chain.Push(l)

	if err := chain.RequireKeys([]string{"A", "B"}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
