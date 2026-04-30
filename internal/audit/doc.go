// Package audit implements a lightweight change-tracking log for envchain.
//
// It records every mutation (set, delete, merge) applied to environment
// config layers, capturing the layer name, key, old value, new value, and
// a timestamp. Entries can be filtered by layer or key and rendered as a
// human-readable summary.
//
// Typical usage:
//
//	log := audit.NewLog(nil) // nil uses time.Now
//	log.Record("dev", "DB_HOST", "", "localhost", audit.ChangeSet)
//	fmt.Print(log.Summary())
//
// The Log is not safe for concurrent use; callers that share a Log across
// goroutines must provide their own synchronisation.
package audit
