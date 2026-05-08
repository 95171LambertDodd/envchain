package compare_test

import (
	"testing"

	"github.com/nicholasgasior/envchain/internal/compare"
)

// stubGetter is a simple in-memory Getter for testing.
type stubGetter struct {
	data map[string]string
}

func newStub(pairs ...string) *stubGetter {
	m := make(map[string]string)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return &stubGetter{data: m}
}

func (s *stubGetter) Keys() []string {
	keys := make([]string, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}
	return keys
}

func (s *stubGetter) Get(key string) (string, bool) {
	v, ok := s.data[key]
	return v, ok
}

func TestNewComparerNilLeft(t *testing.T) {
	_, err := compare.NewComparer(nil, newStub())
	if err == nil {
		t.Fatal("expected error for nil left, got nil")
	}
}

func TestNewComparerNilRight(t *testing.T) {
	_, err := compare.NewComparer(newStub(), nil)
	if err == nil {
		t.Fatal("expected error for nil right, got nil")
	}
}

func TestCompareSameValues(t *testing.T) {
	left := newStub("HOST", "localhost", "PORT", "8080")
	right := newStub("HOST", "localhost", "PORT", "8080")
	cmp, err := compare.NewComparer(left, right)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	res := cmp.Compare()
	if len(res.Same) != 2 {
		t.Errorf("expected 2 same keys, got %d", len(res.Same))
	}
	if len(res.Changed) != 0 || len(res.OnlyLeft) != 0 || len(res.OnlyRight) != 0 {
		t.Errorf("expected no diffs, got changed=%v onlyLeft=%v onlyRight=%v",
			res.Changed, res.OnlyLeft, res.OnlyRight)
	}
}

func TestCompareChangedValues(t *testing.T) {
	left := newStub("HOST", "localhost")
	right := newStub("HOST", "prod.example.com")
	cmp, _ := compare.NewComparer(left, right)
	res := cmp.Compare()
	pair, ok := res.Changed["HOST"]
	if !ok {
		t.Fatal("expected HOST in Changed")
	}
	if pair[0] != "localhost" || pair[1] != "prod.example.com" {
		t.Errorf("unexpected pair: %v", pair)
	}
}

func TestCompareOnlyLeft(t *testing.T) {
	left := newStub("SECRET", "abc", "HOST", "localhost")
	right := newStub("HOST", "localhost")
	cmp, _ := compare.NewComparer(left, right)
	res := cmp.Compare()
	if len(res.OnlyLeft) != 1 || res.OnlyLeft[0] != "SECRET" {
		t.Errorf("expected OnlyLeft=[SECRET], got %v", res.OnlyLeft)
	}
}

func TestCompareOnlyRight(t *testing.T) {
	left := newStub("HOST", "localhost")
	right := newStub("HOST", "localhost", "NEW_KEY", "value")
	cmp, _ := compare.NewComparer(left, right)
	res := cmp.Compare()
	if len(res.OnlyRight) != 1 || res.OnlyRight[0] != "NEW_KEY" {
		t.Errorf("expected OnlyRight=[NEW_KEY], got %v", res.OnlyRight)
	}
}
