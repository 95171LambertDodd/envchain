package tag

import (
	"testing"
)

func makeEntries() map[string]string {
	return map[string]string{
		"DB_PASS":   "hunter2",
		"API_KEY":   "abc123",
		"LOG_LEVEL": "info",
		"REGION":    "us-east-1",
	}
}

func TestFilterByTagNilTaggerError(t *testing.T) {
	_, err := FilterByTag(nil, makeEntries(), "secret")
	if err == nil {
		t.Fatal("expected error for nil tagger")
	}
}

func TestFilterByTagNoTagsError(t *testing.T) {
	tr := NewTagger()
	_, err := FilterByTag(tr, makeEntries())
	if err == nil {
		t.Fatal("expected error when no tags provided")
	}
}

func TestFilterByTagSingleTag(t *testing.T) {
	tr := NewTagger()
	_ = tr.Tag("DB_PASS", "secret")
	_ = tr.Tag("API_KEY", "secret")
	_ = tr.Tag("LOG_LEVEL", "infra")

	keys, err := FilterByTag(tr, makeEntries(), "secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(keys) != 2 || keys[0] != "API_KEY" || keys[1] != "DB_PASS" {
		t.Errorf("unexpected keys: %v", keys)
	}
}

func TestFilterByTagMultipleTags(t *testing.T) {
	tr := NewTagger()
	_ = tr.Tag("DB_PASS", "secret", "sensitive")
	_ = tr.Tag("API_KEY", "secret")

	keys, err := FilterByTag(tr, makeEntries(), "secret", "sensitive")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(keys) != 1 || keys[0] != "DB_PASS" {
		t.Errorf("expected only DB_PASS, got %v", keys)
	}
}

func TestFilterByTagNoMatches(t *testing.T) {
	tr := NewTagger()
	keys, err := FilterByTag(tr, makeEntries(), "nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(keys) != 0 {
		t.Errorf("expected empty result, got %v", keys)
	}
}

func TestGroupByTag(t *testing.T) {
	tr := NewTagger()
	_ = tr.Tag("DB_PASS", "secret", "sensitive")
	_ = tr.Tag("API_KEY", "secret")
	_ = tr.Tag("LOG_LEVEL", "infra")

	groups := GroupByTag(tr, makeEntries())

	if len(groups["secret"]) != 2 {
		t.Errorf("expected 2 keys under 'secret', got %v", groups["secret"])
	}
	if len(groups["sensitive"]) != 1 || groups["sensitive"][0] != "DB_PASS" {
		t.Errorf("expected DB_PASS under 'sensitive', got %v", groups["sensitive"])
	}
	if len(groups["infra"]) != 1 || groups["infra"][0] != "LOG_LEVEL" {
		t.Errorf("expected LOG_LEVEL under 'infra', got %v", groups["infra"])
	}
	if _, ok := groups["REGION"]; ok {
		t.Error("REGION has no tags and should not appear in groups")
	}
}
