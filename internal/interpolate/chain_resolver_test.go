package interpolate_test

import (
	"testing"

	"github.com/yourorg/envchain/internal/interpolate"
)

func TestChainResolverFirstWins(t *testing.T) {
	r1 := interpolate.MapResolver(map[string]string{"KEY": "first"})
	r2 := interpolate.MapResolver(map[string]string{"KEY": "second"})
	r := interpolate.ChainResolver(r1, r2)

	v, ok := r("KEY")
	if !ok {
		t.Fatal("expected key to be found")
	}
	if v != "first" {
		t.Errorf("expected 'first', got %q", v)
	}
}

func TestChainResolverFallsThrough(t *testing.T) {
	r1 := interpolate.MapResolver(map[string]string{})
	r2 := interpolate.MapResolver(map[string]string{"KEY": "fallback"})
	r := interpolate.ChainResolver(r1, r2)

	v, ok := r("KEY")
	if !ok {
		t.Fatal("expected key to be found in second resolver")
	}
	if v != "fallback" {
		t.Errorf("expected 'fallback', got %q", v)
	}
}

func TestChainResolverMissing(t *testing.T) {
	r := interpolate.ChainResolver(
		interpolate.MapResolver(map[string]string{}),
	)
	_, ok := r("MISSING")
	if ok {
		t.Error("expected key to be missing")
	}
}

func TestExpandWithChainResolver(t *testing.T) {
	r := interpolate.ChainResolver(
		interpolate.MapResolver(map[string]string{"APP": "envchain"}),
		interpolate.MapResolver(map[string]string{"ENV": "staging"}),
	)
	i := interpolate.NewInterpolator(r)
	out, err := i.Expand("${APP}-${ENV}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "envchain-staging" {
		t.Errorf("got %q", out)
	}
}
