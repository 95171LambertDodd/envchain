package cache

import (
	"testing"
	"time"
)

func TestNewCacheInvalidCapacity(t *testing.T) {
	_, err := New(0)
	if err != ErrInvalidCapacity {
		t.Fatalf("expected ErrInvalidCapacity, got %v", err)
	}
}

func TestNewCacheValidCapacity(t *testing.T) {
	c, err := New(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Len() != 0 {
		t.Fatalf("expected empty cache")
	}
}

func TestSetAndGet(t *testing.T) {
	c, _ := New(10)
	c.Set("KEY", "value")
	v, err := c.Get("KEY")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "value" {
		t.Fatalf("expected 'value', got %q", v)
	}
}

func TestGetMissingKey(t *testing.T) {
	c, _ := New(10)
	_, err := c.Get("MISSING")
	if err != ErrCacheMiss {
		t.Fatalf("expected ErrCacheMiss, got %v", err)
	}
}

func TestSetWithTTLNotExpired(t *testing.T) {
	c, _ := New(10)
	now := time.Now()
	c.clock = func() time.Time { return now }
	c.SetWithTTL("K", "v", 5*time.Second)
	c.clock = func() time.Time { return now.Add(4 * time.Second) }
	v, err := c.Get("K")
	if err != nil || v != "v" {
		t.Fatalf("expected hit, got err=%v val=%q", err, v)
	}
}

func TestSetWithTTLExpired(t *testing.T) {
	c, _ := New(10)
	now := time.Now()
	c.clock = func() time.Time { return now }
	c.SetWithTTL("K", "v", 5*time.Second)
	c.clock = func() time.Time { return now.Add(6 * time.Second) }
	_, err := c.Get("K")
	if err != ErrCacheMiss {
		t.Fatalf("expected ErrCacheMiss after expiry, got %v", err)
	}
}

func TestDelete(t *testing.T) {
	c, _ := New(10)
	c.Set("K", "v")
	c.Delete("K")
	_, err := c.Get("K")
	if err != ErrCacheMiss {
		t.Fatalf("expected ErrCacheMiss after delete, got %v", err)
	}
}

func TestEvictionAtCapacity(t *testing.T) {
	c, _ := New(2)
	c.Set("A", "1")
	c.Set("B", "2")
	c.Set("C", "3") // should evict one of A or B
	if c.Len() != 2 {
		t.Fatalf("expected 2 items after eviction, got %d", c.Len())
	}
	_, errC := c.Get("C")
	if errC != nil {
		t.Fatalf("newly inserted key C should be present")
	}
}

func TestUpdateExistingKeyDoesNotEvict(t *testing.T) {
	c, _ := New(2)
	c.Set("A", "1")
	c.Set("B", "2")
	c.Set("A", "updated") // update, not new — should not evict
	if c.Len() != 2 {
		t.Fatalf("expected 2 items, got %d", c.Len())
	}
	v, _ := c.Get("A")
	if v != "updated" {
		t.Fatalf("expected 'updated', got %q", v)
	}
}
