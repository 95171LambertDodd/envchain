// Package group provides key-grouping functionality for envchain.
//
// Keys from any Source (a Layer, Chain, Scope, etc.) can be partitioned
// into named buckets using one of three strategies:
//
//   - GroupByPrefix: splits on the first occurrence of a separator,
//     e.g. "DB_HOST" → group "DB" with separator "_".
//
//   - GroupBySuffix: splits on the last occurrence of a separator,
//     e.g. "host_prod" → group "prod" with separator "_".
//
//   - GroupByClassifier: delegates group assignment to a caller-supplied
//     function, enabling arbitrary bucketing logic.
//
// Example usage:
//
//	g, err := group.NewGrouper(chain, group.GroupByPrefix, "_", nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	buckets := g.Group()
//	for name, entries := range buckets {
//		fmt.Printf("group %s: %v\n", name, entries)
//	}
package group
