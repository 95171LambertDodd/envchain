// Package watch provides file-based change detection for environment config layers.
// It monitors a set of file paths and emits reload signals when changes are detected.
package watch

import (
	"os"
	"sync"
	"time"
)

// FileState records the last known modification time of a file.
type FileState struct {
	Path    string
	ModTime time.Time
}

// Watcher polls a set of file paths and notifies via a channel when any file changes.
type Watcher struct {
	mu       sync.Mutex
	paths    []string
	states   map[string]time.Time
	interval time.Duration
	Changes  chan string
	stop     chan struct{}
}

// NewWatcher creates a Watcher that polls the given paths at the given interval.
func NewWatcher(paths []string, interval time.Duration) *Watcher {
	w := &Watcher{
		paths:    paths,
		states:   make(map[string]time.Time),
		interval: interval,
		Changes:  make(chan string, len(paths)),
		stop:     make(chan struct{}),
	}
	for _, p := range paths {
		if info, err := os.Stat(p); err == nil {
			w.states[p] = info.ModTime()
		}
	}
	return w
}

// Start begins polling in a background goroutine.
func (w *Watcher) Start() {
	go func() {
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				w.poll()
			case <-w.stop:
				return
			}
		}
	}()
}

// Stop halts the polling goroutine.
func (w *Watcher) Stop() {
	close(w.stop)
}

func (w *Watcher) poll() {
	w.mu.Lock()
	defer w.mu.Unlock()
	for _, p := range w.paths {
		info, err := os.Stat(p)
		if err != nil {
			continue
		}
		prev, known := w.states[p]
		if !known || info.ModTime().After(prev) {
			w.states[p] = info.ModTime()
			select {
			case w.Changes <- p:
			default:
			}
		}
	}
}
