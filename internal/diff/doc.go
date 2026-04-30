// Package diff computes structured differences between two environment
// Resolvers — any source that maps string keys to string values.
//
// # Overview
//
// A Resolver must implement Get and Keys so the Differ can enumerate
// and compare entries from a base environment against a head environment.
//
// # Change kinds
//
//   - Added   — key present in head but not in base
//   - Removed — key present in base but not in head
//   - Modified — key present in both but with differing values
//
// # Example
//
//	base := myResolver(oldEnv)
//	head := myResolver(newEnv)
//	changes := diff.NewDiffer().Diff(base, head)
//	for _, c := range changes {
//		fmt.Println(c)
//	}
package diff
