package tag

import (
	"testing"
)

func TestTagEmptyKeyError(t *testing.T) {
	tr := NewTagger()
	if err := tr.Tag("", "secret"); err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestTagEmptyTagError(t *testing.T) {
	tr := NewTagger()
	if err := tr.Tag("DB_PASS", ""); err == nil {
		t.Fatal("expected error for empty tag")
	}
}

func TestTagAndRetrieve(t *testing.T) {
	tr := NewTagger()
	if err := tr.Tag("DB_PASS", "secret", "sensitive"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tags := tr.Tags("DB_PASS")
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(tags))
	}
	if tags[0] != "secret" || tags[1] != "sensitive" {
		t.Errorf("unexpected tags: %v", tags)
	}
}

func TestTagsDeduplication(t *testing.T) {
	tr := NewTagger()
	_ = tr.Tag("API_KEY", "secret")
	_ = tr.Tag("API_KEY", "secret") // duplicate
	if len(tr.Tags("API_KEY")) != 1 {
		t.Fatal("expected deduplication of tags")
	}
}

func TestTagsMissingKey(t *testing.T) {
	tr := NewTagger()
	if tr.Tags("NONEXISTENT") != nil {
		t.Fatal("expected nil for unknown key")
	}
}

func TestKeysWithTag(t *testing.T) {
	tr := NewTagger()
	_ = tr.Tag("DB_PASS", "secret")
	_ = tr.Tag("API_KEY", "secret")
	_ = tr.Tag("LOG_LEVEL", "infra")

	keys := tr.KeysWithTag("secret")
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
	if keys[0] != "API_KEY" || keys[1] != "DB_PASS" {
		t.Errorf("unexpected keys: %v", keys)
	}
}

func TestRemoveTag(t *testing.T) {
	tr := NewTagger()
	_ = tr.Tag("DB_PASS", "secret", "sensitive")
	tr.Remove("DB_PASS", "sensitive")
	tags := tr.Tags("DB_PASS")
	if len(tags) != 1 || tags[0] != "secret" {
		t.Errorf("expected only 'secret' after removal, got %v", tags)
	}
}

func TestAllTags(t *testing.T) {
	tr := NewTagger()
	_ = tr.Tag("DB_PASS", "secret")
	_ = tr.Tag("API_KEY", "secret", "infra")
	_ = tr.Tag("LOG_LEVEL", "infra")

	all := tr.AllTags()
	if len(all) != 2 {
		t.Fatalf("expected 2 unique tags, got %d: %v", len(all), all)
	}
	if all[0] != "infra" || all[1] != "secret" {
		t.Errorf("unexpected AllTags order: %v", all)
	}
}

func TestRemoveNoopOnMissingKey(t *testing.T) {
	tr := NewTagger()
	// should not panic
	tr.Remove("NONEXISTENT", "secret")
}
