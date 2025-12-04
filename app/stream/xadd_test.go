package stream

import (
	"strings"
	"testing"
)

func TestXAdd_SingleFieldValue(t *testing.T) {
	store := NewStore()
	result := store.XAdd([]string{"XADD", "stream_key", "0-1", "foo", "bar"})

	// Should return the entry ID as a bulk string
	expected := "$3\r\n0-1\r\n"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}

	// Verify stream was created
	if !store.HasKey("stream_key") {
		t.Error("Expected stream to be created")
	}
}

func TestXAdd_MultipleFieldValues(t *testing.T) {
	store := NewStore()
	result := store.XAdd([]string{"XADD", "stream_key", "1526919030474-0", "temperature", "36", "humidity", "95"})

	// Should return the entry ID as a bulk string
	expected := "$15\r\n1526919030474-0\r\n"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}

	// Verify stream was created
	if !store.HasKey("stream_key") {
		t.Error("Expected stream to be created")
	}
}

func TestXAdd_CreatesStreamIfNotExists(t *testing.T) {
	store := NewStore()

	// Verify stream doesn't exist
	if store.HasKey("newstream") {
		t.Error("Expected stream to not exist initially")
	}

	// Add entry
	store.XAdd([]string{"XADD", "newstream", "0-1", "field", "value"})

	// Verify stream was created
	if !store.HasKey("newstream") {
		t.Error("Expected stream to be created")
	}
}

func TestXAdd_AppendsToExistingStream(t *testing.T) {
	store := NewStore()

	// Add first entry
	store.XAdd([]string{"XADD", "mystream", "0-1", "field1", "value1"})

	// Add second entry
	result := store.XAdd([]string{"XADD", "mystream", "0-2", "field2", "value2"})

	// Should return the second entry ID
	expected := "$3\r\n0-2\r\n"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}

	// Verify stream still exists
	if !store.HasKey("mystream") {
		t.Error("Expected stream to still exist")
	}

	// Verify stream has 2 entries
	store.mutex.Lock()
	entriesCount := len(store.storage["mystream"].entries)
	store.mutex.Unlock()

	if entriesCount != 2 {
		t.Errorf("Expected 2 entries, got %d", entriesCount)
	}
}

func TestXAdd_InvalidArguments(t *testing.T) {
	store := NewStore()

	tests := []struct {
		name string
		args []string
	}{
		{"no arguments", []string{"XADD"}},
		{"only key", []string{"XADD", "key"}},
		{"only key and ID", []string{"XADD", "key", "0-1"}},
		{"missing field value", []string{"XADD", "key", "0-1", "field"}},
		{"odd number of field-value pairs", []string{"XADD", "key", "0-1", "f1", "v1", "f2"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := store.XAdd(tt.args)
			if !strings.HasPrefix(result, "-ERR") {
				t.Errorf("Expected error for %s, got %q", tt.name, result)
			}
		})
	}
}

func TestXAdd_StoresFieldsCorrectly(t *testing.T) {
	store := NewStore()

	// Add entry with multiple fields
	store.XAdd([]string{"XADD", "mystream", "1-0", "temp", "36", "humidity", "95", "location", "room1"})

	// Verify fields are stored correctly
	store.mutex.Lock()
	entry := store.storage["mystream"].entries[0]
	store.mutex.Unlock()

	if entry.ID != "1-0" {
		t.Errorf("Expected ID '1-0', got %q", entry.ID)
	}

	expectedFields := map[string]string{
		"temp":     "36",
		"humidity": "95",
		"location": "room1",
	}

	for key, expectedValue := range expectedFields {
		if actualValue, exists := entry.Fields[key]; !exists {
			t.Errorf("Expected field %q to exist", key)
		} else if actualValue != expectedValue {
			t.Errorf("Expected field %q to have value %q, got %q", key, expectedValue, actualValue)
		}
	}
}
