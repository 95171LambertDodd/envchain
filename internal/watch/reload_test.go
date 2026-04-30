package watch_test

import (
	"errors"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourorg/envchain/internal/watch"
)

func TestReloadManagerCallsHandler(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "app.env")
	if err := os.WriteFile(p, []byte("A=1"), 0o644); err != nil {
		t.Fatal(err)
	}

	var callCount atomic.Int32
	w := watch.NewWatcher([]string{p}, 20*time.Millisecond)
	rm := watch.NewReloadManager(w, func(path string) error {
		callCount.Add(1)
		return nil
	})
	rm.Run()
	defer rm.Stop()

	time.Sleep(40 * time.Millisecond)
	if err := os.WriteFile(p, []byte("A=2"), 0o644); err != nil {
		t.Fatal(err)
	}

	time.Sleep(120 * time.Millisecond)
	if callCount.Load() == 0 {
		t.Error("expected handler to be called at least once")
	}
}

func TestReloadManagerRecordsHandlerErrors(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "bad.env")
	if err := os.WriteFile(p, []byte("X=1"), 0o644); err != nil {
		t.Fatal(err)
	}

	handlerErr := errors.New("parse failed")
	w := watch.NewWatcher([]string{p}, 20*time.Millisecond)
	rm := watch.NewReloadManager(w, func(_ string) error {
		return handlerErr
	})
	rm.Run()
	defer rm.Stop()

	time.Sleep(40 * time.Millisecond)
	if err := os.WriteFile(p, []byte("X=2"), 0o644); err != nil {
		t.Fatal(err)
	}

	time.Sleep(120 * time.Millisecond)
	errs := rm.Errors()
	if len(errs) == 0 {
		t.Fatal("expected at least one recorded error")
	}
	if !errors.Is(errs[0], handlerErr) {
		t.Errorf("unexpected error: %v", errs[0])
	}
}
