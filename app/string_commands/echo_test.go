package string_commands

import (
	"testing"
)

func TestEchoCommand(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "ECHO command with single argument",
			input:    []string{"ECHO", "hello"},
			expected: "$5\r\nhello\r\n",
		},
		{
			name:     "ECHO command with multiple arguments",
			input:    []string{"ECHO", "hello", "world", "test"},
			expected: "$16\r\nhello world test\r\n",
		},
		{
			name:     "ECHO command without arguments",
			input:    []string{"ECHO"},
			expected: "$0\r\n\r\n",
		},
		{
			name:     "Case sensitivity test - echo lowercase",
			input:    []string{"echo", "test"},
			expected: "$4\r\ntest\r\n",
		},
	}

	store := NewStore()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := store.Echo(tt.input)
			if result != tt.expected {
				t.Errorf("Echo(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestEchoEdgeCases(t *testing.T) {
	store := NewStore()

	// Test ECHO with an empty string argument
	result := store.Echo([]string{"ECHO", ""})
	expected := "$0\r\n\r\n"
	if result != expected {
		t.Errorf("Echo([\"ECHO\", \"\"]) = %q, want %q", result, expected)
	}

	// Test ECHO with spaces in arguments
	result = store.Echo([]string{"ECHO", "hello world", "test"})
	expected = "$16\r\nhello world test\r\n"
	if result != expected {
		t.Errorf("Echo([\"ECHO\", \"hello world\", \"test\"]) = %q, want %q", result, expected)
	}
}
