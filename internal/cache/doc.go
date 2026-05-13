// Package cache provides a bounded, thread-safe in-memory cache for
// environment key-value entries used by envchain.
//
// Features:
//
//   - Fixed-capacity store with simple eviction when full
//   - Optional per-entry TTL via SetWithTTL
//   - Safe for concurrent use by multiple goroutines
//
// Basic usage:
//
//	c, err := cache.New(256)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	c.Set("DB_HOST", "localhost")
//	v, err := c.Get("DB_HOST")
//
//	// With expiry:
//	c.SetWithTTL("SESSION_TOKEN", "abc123", 30*time.Second)
//
// ErrCacheMiss is returned by Get when the key is absent or has expired.
package cache
