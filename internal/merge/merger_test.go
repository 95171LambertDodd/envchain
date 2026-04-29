package merge_test

import (
	"testing"

	"github.com/envchain/envchain/internal/config"
	"github.com/envchain/envchain/internal/merge"
)

func makeLayer(t *testing.T, name string, pairs map[string]string) *config.Layer {
	t.Helper()
	l := config.NewLayer(name)
	for k, v := range pairs {
		if err := l.Set(k, v); err != nil {
			t.Fatalf("makeLayer: %v", err)
		}
	}
	return l
}

func TestMergeOverrideStrategy(t *testing.T) {
	base := makeLayer(t, "base", map[string]string{"HOST": "localhost", "PORT": "5432"})
	override := makeLayer(t, "override", map[string]string{"PORT": "6543", "DEBUG": "true"})

	m := merge.NewMerger(merge.StrategyOverride)
	out, err := m.Merge("result", base, override)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertValue(t, out, "HOST", "localhost")
	assertValue(t, out, "PORT", "6543")
	assertValue(t, out, "DEBUG", "true")
}

func TestMergeKeepFirstStrategy(t *testing.T) {
	base := makeLayer(t, "base", map[string]string{"HOST": "localhost", "PORT": "5432"})
	override := makeLayer(t, "override", map[string]string{"PORT": "6543", "DEBUG": "true"})

	m := merge.NewMerger(merge.StrategyKeepFirst)
	out, err := m.Merge("result", base, override)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertValue(t, out, "PORT", "5432") // first value kept
	assertValue(t, out, "DEBUG", "true") // new key added
}

func TestMergeErrorStrategyOnConflict(t *testing.T) {
	base := makeLayer(t, "base", map[string]string{"HOST": "localhost"})
	dup := makeLayer(t, "dup", map[string]string{"HOST": "remotehost"})

	m := merge.NewMerger(merge.StrategyError)
	_, err := m.Merge("result", base, dup)
	if err == nil {
		t.Fatal("expected conflict error, got nil")
	}
}

func TestMergeErrorStrategyNoConflict(t *testing.T) {
	base := makeLayer(t, "base", map[string]string{"HOST": "localhost"})
	extra := makeLayer(t, "extra", map[string]string{"PORT": "8080"})

	m := merge.NewMerger(merge.StrategyError)
	out, err := m.Merge("result", base, extra)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertValue(t, out, "HOST", "localhost")
	assertValue(t, out, "PORT", "8080")
}

func assertValue(t *testing.T, l *config.Layer, key, want string) {
	t.Helper()
	got, ok := l.Get(key)
	if !ok {
		t.Errorf("key %q not found", key)
		return
	}
	if got != want {
		t.Errorf("key %q: got %q, want %q", key, got, want)
	}
}
