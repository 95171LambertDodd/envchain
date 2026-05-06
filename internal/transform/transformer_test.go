package transform_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envchain/internal/transform"
)

func TestNewTransformerNoFuncsError(t *testing.T) {
	_, err := transform.NewTransformer()
	if err == nil {
		t.Fatal("expected error for empty TransformFunc list")
	}
}

func TestUppercaseKeys(t *testing.T) {
	tr, err := transform.NewTransformer(transform.UppercaseKeys())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out, err := tr.Apply(map[string]string{"db_host": "localhost", "port": "5432"})
	if err != nil {
		t.Fatalf("Apply error: %v", err)
	}
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", out["DB_HOST"])
	}
	if out["PORT"] != "5432" {
		t.Errorf("expected PORT=5432, got %q", out["PORT"])
	}
}

func TestTrimSpace(t *testing.T) {
	tr, err := transform.NewTransformer(transform.TrimSpace())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out, err := tr.Apply(map[string]string{"KEY": "  value  "})
	if err != nil {
		t.Fatalf("Apply error: %v", err)
	}
	if out["KEY"] != "value" {
		t.Errorf("expected trimmed value, got %q", out["KEY"])
	}
}

func TestPrefixKeys(t *testing.T) {
	tr, err := transform.NewTransformer(transform.PrefixKeys("APP_"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out, err := tr.Apply(map[string]string{"HOST": "example.com"})
	if err != nil {
		t.Fatalf("Apply error: %v", err)
	}
	if out["APP_HOST"] != "example.com" {
		t.Errorf("expected APP_HOST, got keys: %v", out)
	}
}

func TestStripPrefix(t *testing.T) {
	tr, err := transform.NewTransformer(transform.StripPrefix("APP_"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out, err := tr.Apply(map[string]string{"APP_HOST": "example.com", "PORT": "80"})
	if err != nil {
		t.Fatalf("Apply error: %v", err)
	}
	if out["HOST"] != "example.com" {
		t.Errorf("expected HOST after strip, got %v", out)
	}
	if out["PORT"] != "80" {
		t.Errorf("expected PORT unchanged, got %v", out)
	}
}

func TestChainedTransforms(t *testing.T) {
	tr, err := transform.NewTransformer(
		transform.TrimSpace(),
		transform.UppercaseKeys(),
		transform.PrefixKeys("ENV_"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out, err := tr.Apply(map[string]string{"host": "  prod.example.com  "})
	if err != nil {
		t.Fatalf("Apply error: %v", err)
	}
	if out["ENV_HOST"] != "prod.example.com" {
		t.Errorf("unexpected result: %v", out)
	}
}

func TestTransformFuncError(t *testing.T) {
	errFn := transform.TransformFunc(func(key, value string) (string, string, error) {
		if strings.Contains(value, "bad") {
			return "", "", fmt.Errorf("bad value detected")
		}
		return key, value, nil
	})
	_ = errFn // silence unused; test below uses inline func instead

	tr, _ := transform.NewTransformer(func(k, v string) (string, string, error) {
		if v == "bad" {
			return "", "", fmt.Errorf("rejected")
		}
		return k, v, nil
	})
	_, err := tr.Apply(map[string]string{"K": "bad"})
	if err == nil {
		t.Fatal("expected error from transform func")
	}
}
