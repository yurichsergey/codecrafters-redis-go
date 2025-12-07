package stream

import (
	"testing"
)

// @todo test is strange we should change it
func TestIDToKey(t *testing.T) {
	tests := []struct {
		id1 string
		id2 string
	}{
		{"0-1", "0-2"},
		{"1-0", "1-1"},
		{"99-0", "100-0"},
		{"1526919030474-0", "1526919030474-1"},
		{"1526919030474-0", "1526919030475-0"},
	}

	for _, tt := range tests {
		t.Run(tt.id1+" < "+tt.id2, func(t *testing.T) {
			key1, err := IDToKey(tt.id1)
			if err != nil {
				t.Fatalf("Failed to generate key for %s: %v", tt.id1, err)
			}
			key2, err := IDToKey(tt.id2)
			if err != nil {
				t.Fatalf("Failed to generate key for %s: %v", tt.id2, err)
			}

			if key1 >= key2 {
				t.Errorf("Expected key for %s to be less than key for %s", tt.id1, tt.id2)
			}

			// Verify length
			if len(key1) != 16 {
				t.Errorf("Expected key length 16, got %d", len(key1))
			}
		})
	}
}

func TestIDToKey_Invalid(t *testing.T) {
	_, err := IDToKey("invalid")
	if err == nil {
		t.Error("Expected error for invalid ID")
	}
}
