package main

import "testing"

func TestSetAndGet(t *testing.T) {
	cache := NewLRUCache(1)

	cache.Set("a", "1")

	got, ok := cache.Get("a")
	if !ok {
		t.Fatal("expected key to exist")
	}

	if got != "1" {
		t.Fatalf("expected value %q, got %q", "1", got)
	}
}

func TestGetMissingKey(t *testing.T) {
	cache := NewLRUCache(1)

	_, ok := cache.Get("missing")
	if ok {
		t.Fatalf("expected missing key to return ok=false")
	}
}

func TestSetExistingKeyUpdatesValue(t *testing.T) {
	cache := NewLRUCache(1)

	cache.Set("a", "1")
	cache.Set("a", "2")

	got, ok := cache.Get("a")
	if !ok {
		t.Fatal("expected key to exist")
	}

	if got != "2" {
		t.Fatalf("expected updated value %q, got %q", "2", got)
	}

	if len(cache.Items) != 1 {
		t.Fatalf("expected cache size 1, got %d", len(cache.Items))
	}
}

func TestEvictsLeastRecentlyUsedItem(t *testing.T) {
	cache := NewLRUCache(2)

	cache.Set("a", "1")
	cache.Set("b", "2")
	cache.Set("c", "3")

	if len(cache.Items) != 2 {
		t.Fatalf("expected cache size 2, got %d", len(cache.Items))
	}

	_, ok := cache.Get("a")
	if ok {
		t.Fatal("expected key 'a' to be evicted")
	}

	got, ok := cache.Get("b")
	if ok {
		t.Fatal("expected key 'b' to exist")
	}

	if got != "2" {
		t.Fatalf("expected value %q, got %q", "2", got)
	}

	got, ok = cache.Get("c")
	if !ok {
		t.Fatal("expected key 'c' to exist")
	}

	if got != "3" {
		t.Fatalf("expected value %q, got %q", "3", got)
	}
}

func TestGetMakesItemRecentlyUsed(t *testing.T) {
	cache := NewLRUCache(2)

	cache.Set("a", "1")
	cache.Set("b", "2")

	_, ok := cache.Get("a")
	if !ok {
		t.Fatal("expected key 'a' to exist")
	}

	cache.Set("c", "3")

	_, ok = cache.Get("b")
	if ok {
		t.Fatal("expected key 'b' to be evicted")
	}

	got, ok := cache.Get("a")
	if !ok {
		t.Fatal("expected key 'a' to still exist")
	}

	if got != "1" {
		t.Fatalf("expected value %q, got %q", "1", got)
	}

	got, ok = cache.Get("c")
	if !ok {
		t.Fatal("expected key 'c' to exist")
	}

	if got != "3" {
		t.Fatalf("expected value %q, got %q", "3", got)
	}
}

func TestSetExistingKeyMakesItemRecentlyUsed(t *testing.T) {
	cache := NewLRUCache(2)

	cache.Set("a", "1")
	cache.Set("b", "2")
	cache.Set("a", "updated")
	cache.Set("c", "3")

	_, ok := cache.Get("b")
	if ok {
		t.Fatal("expected key 'b' to be evicted")
	}

	got, ok := cache.Get("a")
	if !ok {
		t.Fatal("expected key 'a' to still exist")
	}

	if got != "updated" {
		t.Fatalf("expected value %q, got %q", "updated", got)
	}
}

func TestLen(t *testing.T) {
	cache := NewLRUCache(2)

	if cache.Len() != 0 {
		t.Fatalf("expected len 0, got %d", cache.Len())
	}

	cache.Set("a", "1")
	if cache.Len() != 1 {
		t.Fatalf("expected len 1, got %d", cache.Len())
	}

	cache.Set("b", "2")
	if cache.Len() != 2 {
		t.Fatalf("expected len 2, got %d", cache.Len())
	}

	cache.Set("c", "3")
	if cache.Len() != 2 {
		t.Fatalf("expected len to stay at capacity 2, got %d", cache.Len())
	}
}

func TestContains(t *testing.T) {
	cache := NewLRUCache(2)

	cache.Set("a", "1")

	if !cache.Contains("a") {
		t.Fatal("expected cache to contain key 'a'")
	}

	if cache.Contains("missing") {
		t.Fatal("expected cache not to contain key 'missing'")
	}
}

func TestDelete(t *testing.T) {
	cache := NewLRUCache(2)

	cache.Set("a", "1")
	cache.Set("b", "2")

	deleted := cache.Delete("a")
	if !deleted {
		t.Fatal("expected Delete to return true")
	}

	if cache.Contains("a") {
		t.Fatal("expected key 'a' to be deleted")
	}

	if cache.Len() != 1 {
		t.Fatalf("expected len 1 after delete, got %d", cache.Len())
	}

	deleted = cache.Delete("missing")
	if deleted {
		t.Fatal("expected Delete on missing key to return false")
	}
}

func TestClear(t *testing.T) {
	cache := NewLRUCache(2)

	cache.Set("a", "1")
	cache.Set("b", "2")

	cache.Clear()

	if cache.Len() != 0 {
		t.Fatalf("expected len 0 after clear, got %d", cache.Len())
	}

	if cache.Contains("a") || cache.Contains("b") {
		t.Fatal("expected cache to be empty after clear")
	}

	if cache.Head != nil {
		t.Fatal("expected Head to be nil after clear")
	}

	if cache.Tail != nil {
		t.Fatal("expected Tail to be nil after clear")
	}
}
