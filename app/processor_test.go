package main

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
			expected: "+hello\r\n",
		},
		{
			name:     "ECHO command with multiple arguments",
			input:    []string{"ECHO", "hello", "world", "test"},
			expected: "+hello world test\r\n",
		},
		{
			name:     "ECHO command without arguments",
			input:    []string{"ECHO"},
			expected: "+\r\n",
		},
		{
			name:     "Unknown command",
			input:    []string{"UNKNOWN"},
			expected: "+PONG\r\n",
		},
		{
			name:     "Unknown command with arguments",
			input:    []string{"GET", "key"},
			expected: "+PONG\r\n",
		},
		{
			name:     "Empty input",
			input:    []string{},
			expected: "",
		},
		{
			name:     "Case sensitivity test - ping lowercase",
			input:    []string{"ping"},
			expected: "+PONG\r\n",
		},
		{
			name:     "Case sensitivity test - echo lowercase",
			input:    []string{"echo", "test"},
			expected: "+PONG\r\n",
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
	if result != "" {
		t.Errorf("defineResponse(nil) = %q, want empty string", result)
	}

	// Test ECHO with an empty string argument
	result = processor.ProcessCommand([]string{"ECHO", ""})
	expected := "+\r\n"
	if result != expected {
		t.Errorf("defineResponse([\"ECHO\", \"\"]) = %q, want %q", result, expected)
	}

	// Test ECHO with spaces in arguments
	result = processor.ProcessCommand([]string{"ECHO", "hello world", "test"})
	expected = "+hello world test\r\n"
	if result != expected {
		t.Errorf("defineResponse([\"ECHO\", \"hello world\", \"test\"]) = %q, want %q", result, expected)
	}
}

func BenchmarkDefineResponse(b *testing.B) {
	testCases := []struct {
		name  string
		input []string
	}{
		{"PING", []string{"PING"}},
		{"ECHO single", []string{"ECHO", "hello"}},
		{"ECHO multiple", []string{"ECHO", "hello", "world", "test"}},
		{"Unknown", []string{"GET", "key"}},
	}

	processor := NewProcessor()
	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				processor.ProcessCommand(tc.input)
			}
		})
	}
}
