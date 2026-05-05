// Package promote provides functionality for promoting environment configuration
// entries from one layer (e.g. staging) to another (e.g. production).
//
// # Overview
//
// A Promoter is created with a Strategy that controls conflict resolution:
//
//   - StrategyOverwrite — destination keys are overwritten by source values.
//   - StrategySkip      — keys that already exist in the destination are left unchanged.
//   - StrategyError     — any conflicting key causes the promotion to abort with an error.
//
// # Example
//
//	promoter := promote.NewPromoter(promote.StrategySkip, "DB_URL", "API_KEY")
//	result, err := promoter.Promote(stagingLayer, prodLayer)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("Promoted:", result.Promoted)
//	fmt.Println("Skipped:",  result.Skipped)
package promote
