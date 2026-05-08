package freeze_test

import (
	"errors"
	"testing"

	"github.com/yourorg/envchain/internal/freeze"
)

func TestSetAndGet(t *testing.T) {
	f := freeze.NewFreezer()
	if err := f.Set("KEY", "value"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, ok := f.Get("KEY")
	if !ok || v != "value" {
		t.Fatalf("expected value=%q ok=true, got %q %v", "value", v, ok)
	}
}

func TestSetEmptyKeyError(t *testing.T) {
	f := freeze.NewFreezer()
	if err := f.Set("", "v"); err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestFreezePreventsMutation(t *testing.T) {
	f := freeze.NewFreezer()
	_ = f.Set("DB_URL", "original")
	if err := f.Freeze("DB_URL"); err != nil {
		t.Fatalf("freeze failed: %v", err)
	}
	err := f.Set("DB_URL", "changed")
	if !errors.Is(err, freeze.ErrFrozen) {
		t.Fatalf("expected ErrFrozen, got %v", err)
	}
	v, _ := f.Get("DB_URL")
	if v != "original" {
		t.Fatalf("value should remain original, got %q", v)
	}
}

func TestFreezeEmptyKeyError(t *testing.T) {
	f := freeze.NewFreezer()
	if err := f.Freeze(""); err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestFreezeAlreadyFrozen(t *testing.T) {
	f := freeze.NewFreezer()
	_ = f.Freeze("KEY")
	err := f.Freeze("KEY")
	if !errors.Is(err, freeze.ErrAlreadyFrozen) {
		t.Fatalf("expected ErrAlreadyFrozen, got %v", err)
	}
}

func TestIsFrozen(t *testing.T) {
	f := freeze.NewFreezer()
	if f.IsFrozen("KEY") {
		t.Fatal("should not be frozen before Freeze call")
	}
	_ = f.Freeze("KEY")
	if !f.IsFrozen("KEY") {
		t.Fatal("should be frozen after Freeze call")
	}
}

func TestFrozenKeys(t *testing.T) {
	f := freeze.NewFreezer()
	_ = f.Freeze("A")
	_ = f.Freeze("B")
	keys := f.FrozenKeys()
	if len(keys) != 2 {
		t.Fatalf("expected 2 frozen keys, got %d", len(keys))
	}
}

func TestGetMissingKey(t *testing.T) {
	f := freeze.NewFreezer()
	_, ok := f.Get("MISSING")
	if ok {
		t.Fatal("expected ok=false for missing key")
	}
}
