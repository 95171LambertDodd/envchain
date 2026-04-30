// Package watch implements lightweight file-change detection for envchain.
//
// Overview
//
// The watch package provides two main types:
//
//   - Watcher: polls one or more file paths at a configurable interval and
//     emits the changed path on a buffered channel (Changes) whenever a
//     file's modification time advances.
//
//   - ReloadManager: wraps a Watcher and dispatches each change notification
//     to a user-supplied ReloadHandler. Handler errors are accumulated and
//     retrievable via Errors().
//
// Usage
//
//		w := watch.NewWatcher([]string{"base.env", "prod.env"}, 500*time.Millisecond)
//		rm := watch.NewReloadManager(w, func(path string) error {
//			// re-parse and reload the layer from path
//			return nil
//		})
//		rm.Run()
//		defer rm.Stop()
//
// The package intentionally avoids OS-level inotify/FSEvents so that it
// works uniformly across Linux, macOS, and Windows without CGO dependencies.
package watch
