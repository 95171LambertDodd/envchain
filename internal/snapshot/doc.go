// Package snapshot provides functionality for capturing, persisting, and
// comparing resolved environment chain states.
//
// A Snapshot records the full set of key-value pairs produced by resolving
// a chain at a specific point in time, along with a human-readable label
// (e.g. "prod", "staging-2024-01-15") and a UTC timestamp.
//
// Typical usage:
//
//	// Capture current resolved state
//	snap := snapshot.NewSnapshot("prod", resolvedEntries)
//
//	// Persist to disk
//	if err := snap.Save("/tmp/prod-snap.json"); err != nil {
//		log.Fatal(err)
//	}
//
//	// Load a previous snapshot
//	old, err := snapshot.Load("/tmp/prod-snap-prev.json")
//
//	// Compare two snapshots to detect drift
//	added, removed, changed := snapshot.Diff(old, snap)
//
Snapshots are stored as indented JSON and are safe to commit to version
control for auditing configuration drift between deployments.
package snapshot
