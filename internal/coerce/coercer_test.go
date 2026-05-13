package coerce_test

import (
	"testing"

	"github.com/user/envchain/internal/coerce"
)

// mapSource is a simple in-memory Source for testing.
type mapSource map[string]string

func (m mapSource) Get(key string) (string, bool) {
	v, ok := m[key]
	return v, ok
}

func TestNewCoercerNilSourceError(t *testing.T) {
	_, err := coerce.NewCoercer(nil)
	if err == nil {
		t.Fatal("expected error for nil source")
	}
}

func TestStringFallback(t *testing.T) {
	c, _ := coerce.NewCoercer(mapSource{})
	if got := c.String("MISSING", "default"); got != "default" {
		t.Fatalf("expected 'default', got %q", got)
	}
}

func TestStringPresent(t *testing.T) {
	c, _ := coerce.NewCoercer(mapSource{"APP": "envchain"})
	if got := c.String("APP", "x"); got != "envchain" {
		t.Fatalf("expected 'envchain', got %q", got)
	}
}

func TestIntFallback(t *testing.T) {
	c, _ := coerce.NewCoercer(mapSource{})
	v, err := c.Int("PORT", 9000)
	if err != nil || v != 9000 {
		t.Fatalf("expected 9000, got %d err %v", v, err)
	}
}

func TestIntParsed(t *testing.T) {
	c, _ := coerce.NewCoercer(mapSource{"PORT": "8080"})
	v, err := c.Int("PORT", 0)
	if err != nil || v != 8080 {
		t.Fatalf("expected 8080, got %d err %v", v, err)
	}
}

func TestIntInvalid(t *testing.T) {
	c, _ := coerce.NewCoercer(mapSource{"PORT": "abc"})
	_, err := c.Int("PORT", 0)
	if err == nil {
		t.Fatal("expected error for non-integer value")
	}
}

func TestBoolTrueVariants(t *testing.T) {
	for _, raw := range []string{"true", "1", "yes", "TRUE", "YES"} {
		c, _ := coerce.NewCoercer(mapSource{"FLAG": raw})
		v, err := c.Bool("FLAG", false)
		if err != nil || !v {
			t.Fatalf("expected true for %q, got %v err %v", raw, v, err)
		}
	}
}

func TestBoolFalseVariants(t *testing.T) {
	for _, raw := range []string{"false", "0", "no", "FALSE"} {
		c, _ := coerce.NewCoercer(mapSource{"FLAG": raw})
		v, err := c.Bool("FLAG", true)
		if err != nil || v {
			t.Fatalf("expected false for %q, got %v err %v", raw, v, err)
		}
	}
}

func TestBoolInvalid(t *testing.T) {
	c, _ := coerce.NewCoercer(mapSource{"FLAG": "maybe"})
	_, err := c.Bool("FLAG", false)
	if err == nil {
		t.Fatal("expected error for invalid bool value")
	}
}

func TestFloatParsed(t *testing.T) {
	c, _ := coerce.NewCoercer(mapSource{"RATIO": "0.75"})
	v, err := c.Float("RATIO", 0)
	if err != nil || v != 0.75 {
		t.Fatalf("expected 0.75, got %f err %v", v, err)
	}
}

func TestFloatFallback(t *testing.T) {
	c, _ := coerce.NewCoercer(mapSource{})
	v, err := c.Float("RATIO", 1.5)
	if err != nil || v != 1.5 {
		t.Fatalf("expected 1.5, got %f err %v", v, err)
	}
}

func TestFloatInvalid(t *testing.T) {
	c, _ := coerce.NewCoercer(mapSource{"RATIO": "not-a-float"})
	_, err := c.Float("RATIO", 0)
	if err == nil {
		t.Fatal("expected error for invalid float value")
	}
}
