package list

import (
	"testing"
)

func TestHasKey_ExistingKey(t *testing.T) {
	store := NewStore()

	// Create a list
	store.RPush([]string{"RPUSH", "mylist", "value1"})

	// Test HasKey
	if !store.HasKey("mylist") {
		t.Error("Expected HasKey to return true for existing list")
	}
}

func TestHasKey_MissingKey(t *testing.T) {
	store := NewStore()

	// Test HasKey on non-existent key
	if store.HasKey("missing_list") {
		t.Error("Expected HasKey to return false for missing list")
	}
}

func TestHasKey_EmptyList(t *testing.T) {
	store := NewStore()

	// Create a list and then pop all elements
	store.RPush([]string{"RPUSH", "mylist", "value1"})
	store.LPop([]string{"LPOP", "mylist"})

	// List should not exist after all elements are removed
	if store.HasKey("mylist") {
		t.Error("Expected HasKey to return false for empty list")
	}
}

func TestHasKey_MultipleKeys(t *testing.T) {
	store := NewStore()

	// Create multiple lists
	store.RPush([]string{"RPUSH", "list1", "value1"})
	store.RPush([]string{"RPUSH", "list2", "value2"})
	store.LPush([]string{"LPUSH", "list3", "value3"})

	// Test all keys exist
	if !store.HasKey("list1") {
		t.Error("Expected HasKey to return true for list1")
	}
	if !store.HasKey("list2") {
		t.Error("Expected HasKey to return true for list2")
	}
	if !store.HasKey("list3") {
		t.Error("Expected HasKey to return true for list3")
	}

	// Test non-existent key
	if store.HasKey("list4") {
		t.Error("Expected HasKey to return false for list4")
	}
}

func TestHasKey_AfterMultiplePushes(t *testing.T) {
	store := NewStore()

	// Create a list with multiple elements
	store.RPush([]string{"RPUSH", "mylist", "value1", "value2", "value3"})

	// Key should exist
	if !store.HasKey("mylist") {
		t.Error("Expected HasKey to return true for list with multiple elements")
	}
}

func TestHasKey_AfterPartialPop(t *testing.T) {
	store := NewStore()

	// Create a list with multiple elements
	store.RPush([]string{"RPUSH", "mylist", "value1", "value2", "value3"})

	// Pop one element
	store.LPop([]string{"LPOP", "mylist"})

	// List should still exist
	if !store.HasKey("mylist") {
		t.Error("Expected HasKey to return true for list after partial pop")
	}
}

func TestHasKey_DifferentListOperations(t *testing.T) {
	store := NewStore()

	// Create list with RPUSH
	store.RPush([]string{"RPUSH", "list1", "value1"})

	// Create list with LPUSH
	store.LPush([]string{"LPUSH", "list2", "value2"})

	// Both should exist
	if !store.HasKey("list1") {
		t.Error("Expected HasKey to return true for list created with RPUSH")
	}
	if !store.HasKey("list2") {
		t.Error("Expected HasKey to return true for list created with LPUSH")
	}
}
