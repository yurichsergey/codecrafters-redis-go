package processor

import (
	"testing"
)

func TestDefineResponse(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "PING command",
			input:    []string{"PING"},
			expected: "+PONG\r\n",
		},
		{
			name:     "PING with additional arguments",
			input:    []string{"PING", "extra", "args"},
			expected: "+PONG\r\n",
		},
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
			name:     "Unknown command",
			input:    []string{"UNKNOWN"},
			expected: "+PONG\r\n",
		},
		{
			name:     "Unknown command with arguments",
			input:    []string{"UNKNOWN", "key"},
			expected: "+PONG\r\n",
		},
		{
			name:     "Empty input",
			input:    []string{},
			expected: "$-1\r\n",
		},
		{
			name:     "Case sensitivity test - ping lowercase",
			input:    []string{"ping"},
			expected: "+PONG\r\n",
		},
		{
			name:     "Case sensitivity test - echo lowercase",
			input:    []string{"echo", "test"},
			expected: "$4\r\ntest\r\n",
		},
	}

	processor := NewProcessor()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processor.ProcessCommand(tt.input)
			if result != tt.expected {
				t.Errorf("defineResponse(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestDefineResponseEdgeCases(t *testing.T) {
	processor := NewProcessor()

	// Test with nil slice
	result := processor.ProcessCommand(nil)
	if result != "$-1\r\n" {
		t.Errorf("defineResponse(nil) = %q, want empty string", result)
	}

	// Test ECHO with an empty string argument
	result = processor.ProcessCommand([]string{"ECHO", ""})
	expected := "$0\r\n\r\n"
	if result != expected {
		t.Errorf("defineResponse([\"ECHO\", \"\"]) = %q, want %q", result, expected)
	}

	// Test ECHO with spaces in arguments
	result = processor.ProcessCommand([]string{"ECHO", "hello world", "test"})
	expected = "$16\r\nhello world test\r\n"
	if result != expected {
		t.Errorf("defineResponse([\"ECHO\", \"hello world\", \"test\"]) = %q, want %q", result, expected)
	}
}
