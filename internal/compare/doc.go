// Package compare provides side-by-side comparison of two environment
// sources (layers, chains, or any Getter implementation).
//
// # Overview
//
// Use NewComparer to create a Comparer from any two Getter values, then
// call Compare to obtain a Result describing:
//
//   - Same      – keys whose values are identical on both sides.
//   - Changed   – keys present on both sides with differing values,
//                 mapped to a [left, right] pair.
//   - OnlyLeft  – keys found only in the left source.
//   - OnlyRight – keys found only in the right source.
//
// # Example
//
//	base, _    := config.NewLayer("base")
//	override, _ := config.NewLayer("override")
//
//	base.Set("HOST", "localhost")
//	override.Set("HOST", "prod.example.com")
//	override.Set("PORT", "443")
//
//	cmp, _ := compare.NewComparer(base, override)
//	result := cmp.Compare()
//	// result.Changed["HOST"] == ["localhost", "prod.example.com"]
//	// result.OnlyRight contains "PORT"
package compare
