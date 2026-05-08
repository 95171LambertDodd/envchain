package pin_test

import (
	"errors"
	"testing"

	"github.com/yourorg/envchain/internal/config"
	"github.com/yourorg/envchain/internal/pin"
)

func makeBaseLayer(t *testing.T, entries map[string]string) *config.Layer {
	t.Helper()
	l, err := config.NewLayer("base")
	if err != nil {
		t.Fatalf("NewLayer: %v", err)
	}
	for k, v := range entries {
		if err := l.Set(k, v); err != nil {
			t.Fatalf("Set %q: %v", k, err)
		}
	}
	return l
}

func TestPinFromLayerNilError(t *testing.T) {
	_, err := pin.PinFromLayer(nil)
	if err == nil {
		t.Fatal("expected error for nil layer")
	}
}

func TestPinFromLayerPinsAllKeys(t *testing.T) {
	l := makeBaseLayer(t, map[string]string{"A": "1", "B": "2"})
	p, err := pin.PinFromLayer(l)
	if err != nil {
		t.Fatalf("PinFromLayer: %v", err)
	}
	if got := len(p.Pinned()); got != 2 {
		t.Fatalf("expected 2 pinned keys, got %d", got)
	}
}

func TestValidateChainNoViolations(t *testing.T) {
	base := makeBaseLayer(t, map[string]string{"APP_ENV": "production"})
	p, _ := pin.PinFromLayer(base)

	chain, _ := config.NewChain()
	_ = chain.Push(base)

	if errs := pin.ValidateChain(p, chain); len(errs) != 0 {
		t.Fatalf("expected no violations, got %v", errs)
	}
}

func TestValidateChainDetectsDrift(t *testing.T) {
	base := makeBaseLayer(t, map[string]string{"APP_ENV": "production"})
	p, _ := pin.PinFromLayer(base)

	override, _ := config.NewLayer("override")
	_ = override.Set("APP_ENV", "staging")

	chain, _ := config.NewChain()
	_ = chain.Push(base)
	_ = chain.Push(override)

	errs := pin.ValidateChain(p, chain)
	if len(errs) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(errs))
	}
	if !errors.Is(errs[0], pin.ErrPinViolation) {
		t.Fatalf("expected ErrPinViolation, got %v", errs[0])
	}
}
