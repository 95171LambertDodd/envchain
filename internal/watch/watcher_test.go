package watch_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/envchain/internal/watch"
)

func writeTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempFile: %v", err)
	}
	return p
}

func TestWatcherDetectsChange(t *testing.T) {
	dir := t.TempDir()
	p := writeTempFile(t, dir, "env.cfg", "KEY=val")

	w := watch.NewWatcher([]string{p}, 20*time.Millisecond)
	w.Start()
	defer w.Stop()

	// Allow one poll to baseline.
	time.Sleep(40 * time.Millisecond)

	// Modify the file.
	if err := os.WriteFile(p, []byte("KEY=changed"), 0o644); err != nil {
		t.Fatal(err)
	}

	select {
	case changed := <-w.Changes:
		if changed != p {
			t.Errorf("expected %q, got %q", p, changed)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timed out waiting for change notification")
	}
}

func TestWatcherNoSpuriousFire(t *testing.T) {
	dir := t.TempDir()
	p := writeTempFile(t, dir, "stable.cfg", "KEY=stable")

	w := watch.NewWatcher([]string{p}, 20*time.Millisecond)
	w.Start()
	defer w.Stop()

	time.Sleep(80 * time.Millisecond)

	select {
	case got := <-w.Changes:
		t.Errorf("unexpected change notification for %q", got)
	default:
		// expected: no change
	}
}

func TestWatcherMissingFileSkipped(t *testing.T) {
	w := watch.NewWatcher([]string{"/nonexistent/path/env.cfg"}, 20*time.Millisecond)
	w.Start()
	defer w.Stop()

	time.Sleep(60 * time.Millisecond)

	select {
	case got := <-w.Changes:
		t.Errorf("unexpected notification for missing file: %q", got)
	default:
	}
}
