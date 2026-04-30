package watch

import (
	"fmt"
	"sync"
)

// ReloadHandler is called when a watched file changes.
// The path argument is the file that triggered the reload.
type ReloadHandler func(path string) error

// ReloadManager binds a Watcher to a ReloadHandler and dispatches events.
type ReloadManager struct {
	watcher *Watcher
	handler ReloadHandler
	errors  []error
	mu      sync.Mutex
}

// NewReloadManager creates a ReloadManager using the provided Watcher and handler.
func NewReloadManager(w *Watcher, h ReloadHandler) *ReloadManager {
	return &ReloadManager{watcher: w, handler: h}
}

// Run starts the watcher and processes change events until Stop is called.
func (r *ReloadManager) Run() {
	r.watcher.Start()
	go func() {
		for path := range r.watcher.Changes {
			if err := r.handler(path); err != nil {
				r.mu.Lock()
				r.errors = append(r.errors, fmt.Errorf("reload %q: %w", path, err))
				r.mu.Unlock()
			}
		}
	}()
}

// Stop halts the underlying watcher.
func (r *ReloadManager) Stop() {
	r.watcher.Stop()
}

// Errors returns a snapshot of all handler errors encountered so far.
func (r *ReloadManager) Errors() []error {
	r.mu.Lock()
	defer r.mu.Unlock()
	copy := make([]error, len(r.errors))
	for i, e := range r.errors {
		copy[i] = e
	}
	return copy
}
