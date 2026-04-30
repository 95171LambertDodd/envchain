package audit_test

import (
	"testing"
	"time"

	"github.com/yourorg/envchain/internal/audit"
)

var fixedTime = time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

func fixedClock() time.Time { return fixedTime }

func TestRecordAndEntries(t *testing.T) {
	l := audit.NewLog(fixedClock)
	l.Record("dev", "DB_HOST", "", "localhost", audit.ChangeSet)
	l.Record("dev", "DB_PORT", "", "5432", audit.ChangeSet)

	entries := l.Entries()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %s", entries[0].Key)
	}
	if entries[0].Timestamp != fixedTime {
		t.Errorf("unexpected timestamp")
	}
}

func TestEntriesReturnsCopy(t *testing.T) {
	l := audit.NewLog(fixedClock)
	l.Record("prod", "SECRET", "", "abc", audit.ChangeSet)

	e1 := l.Entries()
	e1[0].Key = "MUTATED"

	e2 := l.Entries()
	if e2[0].Key == "MUTATED" {
		t.Error("Entries() should return a copy, not a reference")
	}
}

func TestFilterByLayer(t *testing.T) {
	l := audit.NewLog(fixedClock)
	l.Record("dev", "A", "", "1", audit.ChangeSet)
	l.Record("prod", "B", "", "2", audit.ChangeSet)
	l.Record("dev", "C", "old", "", audit.ChangeDelete)

	devEntries := l.FilterByLayer("dev")
	if len(devEntries) != 2 {
		t.Errorf("expected 2 dev entries, got %d", len(devEntries))
	}
}

func TestFilterByKey(t *testing.T) {
	l := audit.NewLog(fixedClock)
	l.Record("dev", "PORT", "", "3000", audit.ChangeSet)
	l.Record("prod", "PORT", "3000", "443", audit.ChangeSet)
	l.Record("dev", "HOST", "", "localhost", audit.ChangeSet)

	portEntries := l.FilterByKey("PORT")
	if len(portEntries) != 2 {
		t.Errorf("expected 2 PORT entries, got %d", len(portEntries))
	}
}

func TestSummaryEmpty(t *testing.T) {
	l := audit.NewLog(nil)
	if l.Summary() != "no audit entries" {
		t.Errorf("expected empty summary message")
	}
}

func TestSummaryNonEmpty(t *testing.T) {
	l := audit.NewLog(fixedClock)
	l.Record("staging", "API_KEY", "", "xyz", audit.ChangeMerge)

	s := l.Summary()
	if s == "no audit entries" {
		t.Error("expected non-empty summary")
	}
	if len(s) == 0 {
		t.Error("summary should not be empty string")
	}
}
