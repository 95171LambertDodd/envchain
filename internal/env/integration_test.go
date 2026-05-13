package env_test

import (
	"testing"

	"github.com/envchain/envchain/internal/config"
	"github.com/envchain/envchain/internal/env"
)

// TestLoaderIntegratesWithChain verifies that a layer produced by Loader can
// be pushed onto a Chain and resolved correctly.
func TestLoaderIntegratesWithChain(t *testing.T) {
	t.Setenv("CHAIN_API_KEY", "secret")
	t.Setenv("CHAIN_LOG_LEVEL", "debug")

	loader := env.NewLoader(env.WithPrefix("CHAIN_"))
	layer, err := loader.Load("env-layer")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	base, _ := config.NewLayer("base")
	_ = base.Set("LOG_LEVEL", "info")  // will be overridden
	_ = base.Set("TIMEOUT", "30s")     // only in base

	chain, err := config.NewChain(base, layer)
	if err != nil {
		t.Fatalf("NewChain: %v", err)
	}

	// env-layer wins for LOG_LEVEL
	if v, err := chain.Resolve("LOG_LEVEL"); err != nil || v != "debug" {
		t.Errorf("LOG_LEVEL: got %q err %v", v, err)
	}
	// base provides TIMEOUT
	if v, err := chain.Resolve("TIMEOUT"); err != nil || v != "30s" {
		t.Errorf("TIMEOUT: got %q err %v", v, err)
	}
	// env-layer provides API_KEY
	if v, err := chain.Resolve("API_KEY"); err != nil || v != "secret" {
		t.Errorf("API_KEY: got %q err %v", v, err)
	}
}

// TestLoaderStrictIntegration checks strict mode end-to-end with a chain.
func TestLoaderStrictIntegration(t *testing.T) {
	t.Setenv("SVC_HOST", "0.0.0.0")
	t.Setenv("SVC_PORT", "9090")

	loader := env.NewLoader(
		env.WithPrefix("SVC_"),
		env.WithStrict("HOST", "PORT"),
	)
	layer, err := loader.Load("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	chain, err := config.NewChain(layer)
	if err != nil {
		t.Fatalf("NewChain: %v", err)
	}

	keys := []string{"HOST", "PORT"}
	expected := []string{"0.0.0.0", "9090"}
	for i, k := range keys {
		v, err := chain.Resolve(k)
		if err != nil {
			t.Errorf("%s: %v", k, err)
			continue
		}
		if v != expected[i] {
			t.Errorf("%s: want %q got %q", k, expected[i], v)
		}
	}
}
