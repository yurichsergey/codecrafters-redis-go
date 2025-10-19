package main

import (
	"strings"
	"testing"
)

func TestLRangeCommand(t *testing.T) {
	tests := []struct {
		name     string
		setup    [][]string // Commands to run before the test
		input    []string
		expected string
	}{
		{
			name:     "LRANGE non-existent list",
			input:    []string{"LRANGE", "nonexistent_list", "0", "1"},
			expected: "*0\r\n",
		},
		{
			name:     "LRANGE with single list element",
			setup:    [][]string{{"RPUSH", "test_list", "a"}},
			input:    []string{"LRANGE", "test_list", "0", "0"},
			expected: "*1\r\n$1\r\na\r\n",
		},
		{
			name:     "LRANGE with multiple elements",
			setup:    [][]string{{"RPUSH", "multi_list", "a", "b", "c", "d", "e"}},
			input:    []string{"LRANGE", "multi_list", "0", "4"},
			expected: "*5\r\n$1\r\na\r\n$1\r\nb\r\n$1\r\nc\r\n$1\r\nd\r\n$1\r\ne\r\n",
		},
		{
			name:     "LRANGE with subset of elements",
			setup:    [][]string{{"RPUSH", "subset_list", "a", "b", "c", "d", "e"}},
			input:    []string{"LRANGE", "subset_list", "1", "3"},
			expected: "*3\r\n$1\r\nb\r\n$1\r\nc\r\n$1\r\nd\r\n",
		},
		{
			name:     "LRANGE with start index greater than list length",
			setup:    [][]string{{"RPUSH", "bounds_list", "a", "b", "c"}},
			input:    []string{"LRANGE", "bounds_list", "5", "6"},
			expected: "*0\r\n",
		},
		{
			name:     "LRANGE with stop index beyond list length",
			setup:    [][]string{{"RPUSH", "beyond_list", "a", "b", "c"}},
			input:    []string{"LRANGE", "beyond_list", "1", "10"},
			expected: "*2\r\n$1\r\nb\r\n$1\r\nc\r\n",
		},
		{
			name:     "LRANGE with start index greater than stop index",
			setup:    [][]string{{"RPUSH", "invalid_range", "a", "b", "c"}},
			input:    []string{"LRANGE", "invalid_range", "2", "1"},
			expected: "*0\r\n",
		},
		{
			name:     "LRANGE with empty string elements",
			setup:    [][]string{{"RPUSH", "empty_elements", "", "text", ""}},
			input:    []string{"LRANGE", "empty_elements", "0", "2"},
			expected: "*3\r\n$0\r\n\r\n$4\r\ntext\r\n$0\r\n\r\n",
		},
		{
			name:     "LRANGE with elements of different lengths",
			setup:    [][]string{{"RPUSH", "mixed_length", "short", "verylongstring", "mid"}},
			input:    []string{"LRANGE", "mixed_length", "0", "2"},
			expected: "*3\r\n$5\r\nshort\r\n$14\r\nverylongstring\r\n$3\r\nmid\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor := NewProcessor()

			// Run setup commands
			for _, setupCmd := range tt.setup {
				processor.ProcessCommand(setupCmd)
			}

			// Run the actual LRANGE test
			result := processor.ProcessCommand(tt.input)
			if result != tt.expected {
				t.Errorf("ProcessCommand(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestLRangeCommandErrors(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "LRANGE without enough arguments",
			input:    []string{"LRANGE", "list"},
			expected: "-ERR wrong number of arguments for 'lrange' command\r\n",
		},
		{
			name:     "LRANGE with too many arguments",
			input:    []string{"LRANGE", "list", "0", "1", "extra"},
			expected: "-ERR wrong number of arguments for 'lrange' command\r\n",
		},
		{
			name:     "LRANGE with non-numeric start index",
			input:    []string{"LRANGE", "list", "abc", "1"},
			expected: "-ERR invalid start index\r\n",
		},
		{
			name:     "LRANGE with non-numeric stop index",
			input:    []string{"LRANGE", "list", "0", "xyz"},
			expected: "-ERR invalid stop index\r\n",
		},
	}

	processor := NewProcessor()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processor.ProcessCommand(tt.input)
			if result != tt.expected {
				t.Errorf("ProcessCommand(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestLRangeCaseInsensitivity(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "LRANGE lowercase",
			input:    []string{"lrange", "list_key", "0", "1"},
			expected: "*0\r\n",
		},
		{
			name:     "LRANGE mixed case",
			input:    []string{"LrAnGe", "list_key", "0", "1"},
			expected: "*0\r\n",
		},
	}

	processor := NewProcessor()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processor.ProcessCommand(tt.input)
			if result != tt.expected {
				t.Errorf("ProcessCommand(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func BenchmarkLRangeCommand(b *testing.B) {
	testCases := []struct {
		name  string
		setup []string
		input []string
	}{
		{
			name:  "LRANGE small list",
			setup: []string{"RPUSH", "small_list", "a", "b", "c", "d", "e"},
			input: []string{"LRANGE", "small_list", "0", "4"},
		},
		{
			name:  "LRANGE medium list",
			setup: []string{"RPUSH", "medium_list", strings.Repeat("a", 100), strings.Repeat("b", 100), strings.Repeat("c", 100)},
			input: []string{"LRANGE", "medium_list", "0", "2"},
		},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			processor := NewProcessor()

			// Setup the list for benchmarking
			processor.ProcessCommand(tc.setup)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				processor.ProcessCommand(tc.input)
			}
		})
	}
}

func TestLRangeCommandNegativeIndexes(t *testing.T) {
	tests := []struct {
		name     string
		setup    [][]string
		input    []string
		expected string
	}{
		{
			name:     "LRANGE with last element using -1",
			setup:    [][]string{{"RPUSH", "last_test", "a", "b", "c", "d", "e"}},
			input:    []string{"LRANGE", "last_test", "-1", "-1"},
			expected: "*1\r\n$1\r\ne\r\n",
		},
		{
			name:     "LRANGE with last two elements using -2 -1",
			setup:    [][]string{{"RPUSH", "last_two_test", "a", "b", "c", "d", "e"}},
			input:    []string{"LRANGE", "last_two_test", "-2", "-1"},
			expected: "*2\r\n$1\r\nd\r\n$1\r\ne\r\n",
		},
		{
			name:     "LRANGE with start negative index and end positive index",
			setup:    [][]string{{"RPUSH", "mixed_index", "a", "b", "c", "d", "e"}},
			input:    []string{"LRANGE", "mixed_index", "-3", "4"},
			expected: "*3\r\n$1\r\nc\r\n$1\r\nd\r\n$1\r\ne\r\n",
		},
		{
			name:     "LRANGE with both negative indexes spanning list",
			setup:    [][]string{{"RPUSH", "negative_span", "a", "b", "c", "d", "e"}},
			input:    []string{"LRANGE", "negative_span", "-5", "-1"},
			expected: "*5\r\n$1\r\na\r\n$1\r\nb\r\n$1\r\nc\r\n$1\r\nd\r\n$1\r\ne\r\n",
		},
		{
			name:     "LRANGE with negative index out of range treated as 0",
			setup:    [][]string{{"RPUSH", "out_of_range", "a", "b", "c"}},
			input:    []string{"LRANGE", "out_of_range", "-6", "1"},
			expected: "*2\r\n$1\r\na\r\n$1\r\nb\r\n",
		},
		{
			name:     "LRANGE with negative start index greater than list length",
			setup:    [][]string{{"RPUSH", "too_negative", "a", "b", "c"}},
			input:    []string{"LRANGE", "too_negative", "-10", "-5"},
			expected: "*1\r\n$1\r\na\r\n",
		},
		{
			name:     "LRANGE with negative indexes in reverse order",
			setup:    [][]string{{"RPUSH", "reverse_order", "a", "b", "c", "d", "e"}},
			input:    []string{"LRANGE", "reverse_order", "-1", "-5"},
			expected: "*0\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor := NewProcessor()

			// Run setup commands
			for _, setupCmd := range tt.setup {
				processor.ProcessCommand(setupCmd)
			}

			// Run the actual LRANGE test
			result := processor.ProcessCommand(tt.input)
			if result != tt.expected {
				t.Errorf("ProcessCommand(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
