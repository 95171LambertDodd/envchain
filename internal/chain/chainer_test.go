package chain_test

import (
	"testing"

	"envchain/internal/chain"
)

// mapSource is a trivial in-memory Source used in tests.
type mapSource map[string]string

func (m mapSource) Keys() []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func (m mapSource) Get(key string) (string, bool) {
	v, ok := m[key]
	return v, ok
}

func TestNewChainerRequiresSource(t *testing.T) {
	_, err := chain.New()
	if err == nil {
		t.Fatal("expected error when no sources provided")
	}
}

func TestGetFirstSourceWins(t *testing.T) {
	a := mapSource{"KEY": "from-a"}
	b := mapSource{"KEY": "from-b"}
	c, _ := chain.New(a, b)

	v, ok := c.Get("KEY")
	if !ok {
		t.Fatal("expected key to be found")
	}
	if v != "from-a" {
		t.Fatalf("expected 'from-a', got %q", v)
	}
}

func TestGetFallsThrough(t *testing.T) {
	a := mapSource{"OTHER": "x"}
	b := mapSource{"KEY": "from-b"}
	c, _ := chain.New(a, b)

	v, ok := c.Get("KEY")
	if !ok {
		t.Fatal("expected fallthrough to second source")
	}
	if v != "from-b" {
		t.Fatalf("expected 'from-b', got %q", v)
	}
}

func TestGetMissingKey(t *testing.T) {
	a := mapSource{"A": "1"}
	c, _ := chain.New(a)

	_, ok := c.Get("MISSING")
	if ok {
		t.Fatal("expected false for missing key")
	}
}

func TestKeysDeduplication(t *testing.T) {
	a := mapSource{"SHARED": "a", "ONLY_A": "1"}
	b := mapSource{"SHARED": "b", "ONLY_B": "2"}
	c, _ := chain.New(a, b)

	keys := c.Keys()
	seen := make(map[string]int)
	for _, k := range keys {
		seen[k]++
	}
	if seen["SHARED"] != 1 {
		t.Fatalf("expected SHARED once, got %d", seen["SHARED"])
	}
	if len(keys) != 3 {
		t.Fatalf("expected 3 unique keys, got %d", len(keys))
	}
}

func TestOriginWithLabel(t *testing.T) {
	a := mapSource{"A": "1"}
	b := mapSource{"B": "2"}
	c, _ := chain.New(a, b)
	_ = c.WithLabel(0, "base")
	_ = c.WithLabel(1, "override")

	if got := c.Origin("A"); got != "base" {
		t.Fatalf("expected 'base', got %q", got)
	}
	if got := c.Origin("B"); got != "override" {
		t.Fatalf("expected 'override', got %q", got)
	}
	if got := c.Origin("MISSING"); got != "" {
		t.Fatalf("expected empty origin for missing key, got %q", got)
	}
}

func TestWithLabelOutOfRange(t *testing.T) {
	a := mapSource{"A": "1"}
	c, _ := chain.New(a)

	if err := c.WithLabel(5, "bad"); err == nil {
		t.Fatal("expected error for out-of-range index")
	}
}
