package stream

import (
	"testing"
)

func TestHasKey(t *testing.T) {
	store := NewStore()

	// Test non-existent stream
	if store.HasKey("nonexistent") {
		t.Error("Expected HasKey to return false for non-existent stream")
	}

	// Add a stream entry
	store.XAdd([]string{"XADD", "mystream", "0-1", "field", "value"})

	// Test existing stream
	if !store.HasKey("mystream") {
		t.Error("Expected HasKey to return true for existing stream")
	}

	// Test different non-existent stream
	if store.HasKey("anotherstream") {
		t.Error("Expected HasKey to return false for different non-existent stream")
	}
}
