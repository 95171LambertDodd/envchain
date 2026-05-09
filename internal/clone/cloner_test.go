package clone_test

import (
	"testing"

	"github.com/yourorg/envchain/internal/clone"
	"github.com/yourorg/envchain/internal/config"
)

func makeLayer(t *testing.T, name string, pairs map[string]string) *config.Layer {
	t.Helper()
	l, err := config.NewLayer(name)
	if err != nil {
		t.Fatalf("NewLayer(%q): %v", name, err)
	}
	for k, v := range pairs {
		if err := l.Set(k, v); err != nil {
			t.Fatalf("Set(%q): %v", k, err)
		}
	}
	return l
}

func TestNewClonerEmptySuffixError(t *testing.T) {
	_, err := clone.NewCloner("")
	if err == nil {
		t.Fatal("expected error for empty suffix")
	}
}

func TestCloneLayerCopiesValues(t *testing.T) {
	src := makeLayer(t, "base", map[string]string{"HOST": "localhost", "PORT": "5432"})
	c, _ := clone.NewCloner("-copy")
	dst, err := c.Layer(src)
	if err != nil {
		t.Fatalf("Layer: %v", err)
	}
	if dst.Name() != "base-copy" {
		t.Errorf("name = %q, want %q", dst.Name(), "base-copy")
	}
	v, _ := dst.Get("HOST")
	if v != "localhost" {
		t.Errorf("HOST = %q, want %q", v, "localhost")
	}
}

func TestCloneLayerIsolation(t *testing.T) {
	src := makeLayer(t, "base", map[string]string{"KEY": "original"})
	c, _ := clone.NewCloner("-clone")
	dst, _ := c.Layer(src)
	_ = dst.Set("KEY", "mutated")
	v, _ := src.Get("KEY")
	if v != "original" {
		t.Errorf("src mutated: KEY = %q, want %q", v, "original")
	}
}

func TestCloneLayerNilError(t *testing.T) {
	c, _ := clone.NewCloner("-x")
	_, err := c.Layer(nil)
	if err == nil {
		t.Fatal("expected error for nil layer")
	}
}

func TestCloneChainPreservesOrder(t *testing.T) {
	l1 := makeLayer(t, "base", map[string]string{"A": "1"})
	l2 := makeLayer(t, "override", map[string]string{"A": "2"})
	chain := config.NewChain()
	chain.Push(l1)
	chain.Push(l2)

	c, _ := clone.NewCloner("-dup")
	dup, err := c.Chain(chain)
	if err != nil {
		t.Fatalf("Chain: %v", err)
	}
	layers := dup.Layers()
	if len(layers) != 2 {
		t.Fatalf("layer count = %d, want 2", len(layers))
	}
	if layers[0].Name() != "base-dup" || layers[1].Name() != "override-dup" {
		t.Errorf("unexpected layer names: %q %q", layers[0].Name(), layers[1].Name())
	}
}

func TestCloneChainNilError(t *testing.T) {
	c, _ := clone.NewCloner("-x")
	_, err := c.Chain(nil)
	if err == nil {
		t.Fatal("expected error for nil chain")
	}
}
