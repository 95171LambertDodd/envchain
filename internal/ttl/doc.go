// Package ttl implements time-to-live expiry for individual environment config
// keys. It wraps any Source (a type providing Keys/Get) and overlays a
// per-key deadline map. Once a key's deadline has passed, Get treats it as
// absent and Keys omits it from the listing.
//
// Typical usage:
//
//	store, err := ttl.NewTTLStore(myLayer, time.Now)
//	if err != nil { ... }
//
//	// Expire a secret after 30 minutes.
//	_ = store.SetTTL("DB_PASSWORD", 30*time.Minute)
//
//	// Or set an absolute deadline.
//	_ = store.SetExpiry("API_KEY", rotationDeadline)
//
//	// Reads respect the deadline automatically.
//	val, ok := store.Get("DB_PASSWORD")
//
// The clock function (now) is injectable so tests can control time without
// sleeping. Pass time.Now in production code.
package ttl
