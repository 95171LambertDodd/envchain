// Package redact provides a Redactor type for masking sensitive environment
// variable values in key/value maps.
//
// A Redactor is configured with a list of substrings that identify sensitive
// keys (e.g. "secret", "token", "password") and a RedactMode that controls
// how matching values are obscured:
//
//   - ModeStars:   replaces the entire value with "********"
//   - ModePartial: reveals the first two characters and masks the rest
//   - ModeHash:    replaces the value with a length-hinted placeholder
//
// Example:
//
//	r := redact.NewRedactor([]string{"secret", "token"}, redact.ModePartial)
//	safe := r.RedactMap(myEnvMap)
//
// RedactMap never mutates the input map and is safe to call concurrently
// provided the underlying map is not written to during the call.
package redact
