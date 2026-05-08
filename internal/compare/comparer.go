// Package compare provides utilities for comparing two environment layers
// or chains, producing a structured result of matching, differing, and
// exclusive keys.
package compare

import "fmt"

// Result holds the outcome of comparing two sets of environment entries.
type Result struct {
	// Same contains keys whose values are identical in both sides.
	Same []string
	// Changed maps a key to [leftVal, rightVal] where values differ.
	Changed map[string][2]string
	// OnlyLeft contains keys present only in the left side.
	OnlyLeft []string
	// OnlyRight contains keys present only in the right side.
	OnlyRight []string
}

// Getter is the minimal interface required by the Comparer.
type Getter interface {
	Keys() []string
	Get(key string) (string, bool)
}

// Comparer compares two Getter sources.
type Comparer struct {
	left  Getter
	right Getter
}

// NewComparer creates a Comparer for the given left and right sources.
// Returns an error if either source is nil.
func NewComparer(left, right Getter) (*Comparer, error) {
	if left == nil {
		return nil, fmt.Errorf("compare: left source must not be nil")
	}
	if right == nil {
		return nil, fmt.Errorf("compare: right source must not be nil")
	}
	return &Comparer{left: left, right: right}, nil
}

// Compare performs the comparison and returns a Result.
func (c *Comparer) Compare() Result {
	res := Result{
		Changed:  make(map[string][2]string),
	}

	rightKeys := make(map[string]struct{})
	for _, k := range c.right.Keys() {
		rightKeys[k] = struct{}{}
	}

	for _, k := range c.left.Keys() {
		lv, _ := c.left.Get(k)
		rv, ok := c.right.Get(k)
		if !ok {
			res.OnlyLeft = append(res.OnlyLeft, k)
			continue
		}
		delete(rightKeys, k)
		if lv == rv {
			res.Same = append(res.Same, k)
		} else {
			res.Changed[k] = [2]string{lv, rv}
		}
	}

	for k := range rightKeys {
		res.OnlyRight = append(res.OnlyRight, k)
	}
	return res
}
