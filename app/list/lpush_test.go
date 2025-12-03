package list_test

import (
	"testing"

	"github.com/codecrafters-io/redis-starter-go/app/processor"
)

func TestLPushCommand(t *testing.T) {
	tests := []struct {
		name         string
		setup        [][]string // Commands to run before the test
		input        []string
		expected     string
		expectedList []string
	}{
		{
			name:         "LPUSH to a non-existent list",
			input:        []string{"LPUSH", "new_list", "a"},
			expected:     ":1\r\n",
			expectedList: []string{"a"},
		},
		{
			name:         "LPUSH single element to an existing list",
			setup:        [][]string{{"RPUSH", "test_list", "b", "c"}},
			input:        []string{"LPUSH", "test_list", "a"},
			expected:     ":3\r\n",
			expectedList: []string{"a", "b", "c"},
		},
		{
			name:         "LPUSH multiple elements",
			setup:        [][]string{{"RPUSH", "multi_list", "c", "d"}},
			input:        []string{"LPUSH", "multi_list", "b", "a"},
			expected:     ":4\r\n",
			expectedList: []string{"a", "b", "c", "d"},
		},
		{
			name:         "LPUSH with many elements",
			input:        []string{"LPUSH", "many_list", "c", "b", "a"},
			expected:     ":3\r\n",
			expectedList: []string{"a", "b", "c"},
		},
		{
			name:         "LPUSH with empty strings",
			input:        []string{"LPUSH", "empty_list", "", "text", ""},
			expected:     ":3\r\n",
			expectedList: []string{"", "text", ""},
		},
		{
			name:         "LPUSH with long strings",
			input:        []string{"LPUSH", "long_string_list", "verylongstring", "short"},
			expected:     ":2\r\n",
			expectedList: []string{"short", "verylongstring"},
		},
		{
			name: "LPUSH multiple elements",
			setup: [][]string{
				{"LPUSH", "pear", "raspberry"},
				{"LPUSH", "pear", "pear", "pineapple"},
			},
			input:        []string{"LRANGE", "pear", "0", "-1"},
			expected:     "*3\r\n$9\r\npineapple\r\n$4\r\npear\r\n$9\r\nraspberry\r\n",
			expectedList: []string{"pineapple", "pear", "raspberry"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor := processor.NewProcessor()

			// Run setup commands
			for _, setupCmd := range tt.setup {
				processor.ProcessCommand(setupCmd)
			}

			// Run the actual LPUSH test
			result := processor.ProcessCommand(tt.input)
			if result != tt.expected {
				t.Errorf("ProcessCommand(%v) = %q, want %q", tt.input, result, tt.expected)
			}

			// Verify the order using LRANGE
			verifyCmd := []string{"LRANGE", tt.input[1], "0", "-1"}
			verifyResult := processor.ProcessCommand(verifyCmd)

			// Construct expected LRANGE result based on all elements
			expectedListResponse := constructListResponse(tt.expectedList)
			if verifyResult != expectedListResponse {
				t.Errorf("List order incorrect after LPUSH. Got %q, want %q", verifyResult, expectedListResponse)
			}
		})
	}
}

func TestLPushCommandErrors(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "LPUSH without key",
			input:    []string{"LPUSH"},
			expected: "-ERR wrong number of arguments for 'lpush' command\r\n",
		},
		{
			name:     "LPUSH with only key",
			input:    []string{"LPUSH", "list_key"},
			expected: "-ERR wrong number of arguments for 'lpush' command\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor := processor.NewProcessor()
			result := processor.ProcessCommand(tt.input)
			if result != tt.expected {
				t.Errorf("ProcessCommand(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestLPushCaseInsensitivity(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "LPUSH lowercase",
			input:    []string{"lpush", "list_key", "a"},
			expected: ":1\r\n",
		},
		{
			name:     "LPUSH mixed case",
			input:    []string{"LpUsH", "list_key", "a", "b"},
			expected: ":2\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor := processor.NewProcessor()
			result := processor.ProcessCommand(tt.input)
			if result != tt.expected {
				t.Errorf("ProcessCommand(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func BenchmarkLPushCommand(b *testing.B) {
	testCases := []struct {
		name  string
		input []string
	}{
		{
			name:  "LPUSH single element",
			input: []string{"LPUSH", "benchmark_list", "element"},
		},
		{
			name:  "LPUSH multiple elements",
			input: []string{"LPUSH", "benchmark_list", "element1", "element2", "element3"},
		},
		{
			name:  "LPUSH with long string",
			input: []string{"LPUSH", "benchmark_list", "very_long_string_element_for_benchmarking"},
		},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			processor := processor.NewProcessor()

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				processor.ProcessCommand(tc.input)
			}
		})
	}
}
