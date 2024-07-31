package memorydb

import (
	"testing"
	"time"
)

func TestMemoryDB_SetAndGet(t *testing.T) {
	store := NewMemoryDB()
	store.Set("key1", "value1", 10)
	store.Set("key2", "value2", 0) // 永不过期

	val, err := store.Get("key1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if val != "value1" {
		t.Fatalf("Expected value1, got %v", val)
	}

	val, err = store.Get("key2")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if val != "value2" {
		t.Fatalf("Expected value2, got %v", val)
	}
}

func TestMemoryDB_Expire(t *testing.T) {
	store := NewMemoryDB()
	store.Set("key1", "value1", 0)
	store.Expire("key1", 1)

	time.Sleep(2 * time.Second)

	_, err := store.Get("key1")
	if err != ErrKeyExpired {
		t.Fatalf("Expected ErrKeyExpired, got %v", err)
	}
}

func TestMemoryDB_Del(t *testing.T) {
	store := NewMemoryDB()
	store.Set("key1", "value1", 0)
	store.Del("key1")

	_, err := store.Get("key1")
	if err != ErrKeyNotFound {
		t.Fatalf("Expected ErrKeyNotFound, got %v", err)
	}
}

func TestMemoryDB_Cleanup(t *testing.T) {
	store := NewMemoryDB()
	store.Set("key1", "value1", 1)

	time.Sleep(2 * time.Second)

	_, err := store.Get("key1")
	if err != ErrKeyNotFound && err != ErrKeyExpired {
		t.Fatalf("Expected key to be expired or not found, got %v", err)
	}
}
