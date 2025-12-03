package string_commands

import (
	"testing"
	"time"
)

func TestHasKey_ExistingKey(t *testing.T) {
	store := NewStore()

	// Set a key
	store.Set([]string{"SET", "mykey", "myvalue"})

	// Test HasKey
	if !store.HasKey("mykey") {
		t.Error("Expected HasKey to return true for existing key")
	}
}

func TestHasKey_MissingKey(t *testing.T) {
	store := NewStore()

	// Test HasKey on non-existent key
	if store.HasKey("missing_key") {
		t.Error("Expected HasKey to return false for missing key")
	}
}

func TestHasKey_ExpiredKey(t *testing.T) {
	store := NewStore()

	// Set a key with 50ms expiry
	store.Set([]string{"SET", "expiring_key", "value", "PX", "50"})

	// Key should exist immediately
	if !store.HasKey("expiring_key") {
		t.Error("Expected HasKey to return true for key before expiration")
	}

	// Wait for key to expire
	time.Sleep(100 * time.Millisecond)

	// Key should not exist after expiration
	if store.HasKey("expiring_key") {
		t.Error("Expected HasKey to return false for expired key")
	}
}

func TestHasKey_KeyWithoutExpiry(t *testing.T) {
	store := NewStore()

	// Set a key without expiry
	store.Set([]string{"SET", "permanent_key", "value"})

	// Wait a bit
	time.Sleep(50 * time.Millisecond)

	// Key should still exist
	if !store.HasKey("permanent_key") {
		t.Error("Expected HasKey to return true for key without expiry")
	}
}

func TestHasKey_MultipleKeys(t *testing.T) {
	store := NewStore()

	// Set multiple keys
	store.Set([]string{"SET", "key1", "value1"})
	store.Set([]string{"SET", "key2", "value2"})
	store.Set([]string{"SET", "key3", "value3"})

	// Test all keys exist
	if !store.HasKey("key1") {
		t.Error("Expected HasKey to return true for key1")
	}
	if !store.HasKey("key2") {
		t.Error("Expected HasKey to return true for key2")
	}
	if !store.HasKey("key3") {
		t.Error("Expected HasKey to return true for key3")
	}

	// Test non-existent key
	if store.HasKey("key4") {
		t.Error("Expected HasKey to return false for key4")
	}
}

func TestHasKey_OverwrittenKey(t *testing.T) {
	store := NewStore()

	// Set a key with expiry
	store.Set([]string{"SET", "mykey", "value1", "PX", "50"})

	// Overwrite with new value and no expiry
	store.Set([]string{"SET", "mykey", "value2"})

	// Wait past the original expiry time
	time.Sleep(100 * time.Millisecond)

	// Key should still exist because it was overwritten without expiry
	if !store.HasKey("mykey") {
		t.Error("Expected HasKey to return true for overwritten key without expiry")
	}
}
