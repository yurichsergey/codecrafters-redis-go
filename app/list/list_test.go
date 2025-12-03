package list_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/codecrafters-io/redis-starter-go/app/processor"
)

func TestRPushCommand(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "RPUSH creating a new list with a single element",
			input:    []string{"RPUSH", "list_key", "foo"},
			expected: ":1\r\n",
		},
		{
			name:     "RPUSH adding another element to the same list",
			input:    []string{"RPUSH", "list_key", "bar"},
			expected: ":2\r\n",
		},
		{
			name:     "RPUSH with multiple elements in one call",
			input:    []string{"RPUSH", "another_list", "a", "b", "c"},
			expected: ":3\r\n",
		},
		{
			name:     "RPUSH with empty string as element",
			input:    []string{"RPUSH", "empty_list", ""},
			expected: ":1\r\n",
		},
		{
			name:     "RPUSH without enough arguments",
			input:    []string{"RPUSH", "list_key"},
			expected: "-ERR wrong number of arguments for 'rpush' command\r\n",
		},
		{
			name:     "RPUSH case-insensitive command",
			input:    []string{"rpush", "case_list", "element"},
			expected: ":1\r\n",
		},
	}

	processor := processor.NewProcessor()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processor.ProcessCommand(tt.input)
			if result != tt.expected {
				t.Errorf("ProcessCommand(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

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
			processor := processor.NewProcessor()

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

	processor := processor.NewProcessor()
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

	processor := processor.NewProcessor()
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
			processor := processor.NewProcessor()

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
			processor := processor.NewProcessor()

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

func constructListResponse(values []string) string {
	// Construct the RESP array response
	var response string
	response = fmt.Sprintf("*%d\r\n", len(values))
	for _, item := range values {
		response += fmt.Sprintf("$%d\r\n%s\r\n", len(item), item)
	}

	return response
}

func TestLLenCommand(t *testing.T) {
	tests := []struct {
		name     string
		setup    [][]string // Commands to run before the test
		input    []string   // LLEN command
		expected string     // Expected response
	}{
		{
			name:     "LLEN on empty list",
			setup:    nil,
			input:    []string{"LLEN", "non_existent_list"},
			expected: ":0\r\n",
		},
		{
			name:     "LLEN on single item list",
			setup:    [][]string{{"RPUSH", "test_list", "item1"}},
			input:    []string{"LLEN", "test_list"},
			expected: ":1\r\n",
		},
		{
			name:     "LLEN on multiple item list",
			setup:    [][]string{{"RPUSH", "test_list", "item1", "item2", "item3", "item4", "item5"}},
			input:    []string{"LLEN", "test_list"},
			expected: ":5\r\n",
		},
		{
			name:     "LLEN after adding items with multiple commands",
			setup:    [][]string{{"RPUSH", "multi_list", "a"}, {"LPUSH", "multi_list", "b"}, {"RPUSH", "multi_list", "c"}},
			input:    []string{"LLEN", "multi_list"},
			expected: ":3\r\n",
		},
		{
			name:     "LLEN missing arguments",
			setup:    nil,
			input:    []string{"LLEN"},
			expected: "-ERR wrong number of arguments for 'llen' command\r\n",
		},
		{
			name:     "LLEN with too many arguments",
			setup:    nil,
			input:    []string{"LLEN", "list_key", "extra_arg"},
			expected: "-ERR wrong number of arguments for 'llen' command\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor := processor.NewProcessor()

			// Run setup commands
			for _, setupCmd := range tt.setup {
				processor.ProcessCommand(setupCmd)
			}

			// Run the actual LLEN command
			result := processor.ProcessCommand(tt.input)
			if result != tt.expected {
				t.Errorf("ProcessCommand(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestLPopCommand(t *testing.T) {
	tests := []struct {
		name         string
		setup        [][]string // Commands to run before the test
		input        []string
		expected     string
		expectedList []string
	}{
		{
			name:     "LPOP on non-existent list",
			input:    []string{"LPOP", "nonexistent_list"},
			expected: "$-1\r\n",
		},
		{
			name:         "LPOP from a single-element list",
			setup:        [][]string{{"RPUSH", "single_list", "a"}},
			input:        []string{"LPOP", "single_list"},
			expected:     "$1\r\na\r\n",
			expectedList: []string{},
		},
		{
			name:         "LPOP from a multi-element list",
			setup:        [][]string{{"RPUSH", "multi_list", "a", "b", "c", "d"}},
			input:        []string{"LPOP", "multi_list"},
			expected:     "$1\r\na\r\n",
			expectedList: []string{"b", "c", "d"},
		},
		{
			name:         "LPOP multiple times from a list",
			setup:        [][]string{{"RPUSH", "multiple_list", "a", "b", "c", "d"}},
			input:        []string{"LPOP", "multiple_list"},
			expected:     "$1\r\na\r\n",
			expectedList: []string{"b", "c", "d"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor := processor.NewProcessor()

			// Run setup commands
			for _, setupCmd := range tt.setup {
				processor.ProcessCommand(setupCmd)
			}

			// Run the actual LPOP test
			result := processor.ProcessCommand(tt.input)
			if result != tt.expected {
				t.Errorf("ProcessCommand(%v) = %q, want %q", tt.input, result, tt.expected)
			}

			// Verify the order using LRANGE
			verifyCmd := []string{"LRANGE", tt.input[1], "0", "-1"}
			verifyResult := processor.ProcessCommand(verifyCmd)

			// Construct expected LRANGE result
			expectedListResponse := constructListResponse(tt.expectedList)
			if verifyResult != expectedListResponse {
				t.Errorf("List order incorrect after LPOP. Got %q, want %q", verifyResult, expectedListResponse)
			}
		})
	}
}

func TestLPopCommandErrors(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "LPOP without key",
			input:    []string{"LPOP"},
			expected: "-ERR wrong number of arguments for 'lpop' command\r\n",
		},
		{
			name:     "LPOP with too many arguments",
			input:    []string{"LPOP", "list_key", "extra"},
			expected: "-ERR value is not an integer or out of range\r\n",
		},
	}

	processor := processor.NewProcessor()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processor.ProcessCommand(tt.input)
			if result != tt.expected {
				t.Errorf("ProcessCommand(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestLPopCaseInsensitivity(t *testing.T) {
	tests := []struct {
		name     string
		setup    [][]string
		input    []string
		expected string
	}{
		{
			name:     "LPOP lowercase",
			setup:    [][]string{{"RPUSH", "list_key", "a", "b", "c"}},
			input:    []string{"lpop", "list_key"},
			expected: "$1\r\na\r\n",
		},
		{
			name:     "LPOP mixed case",
			setup:    [][]string{{"RPUSH", "list_key", "a", "b", "c"}},
			input:    []string{"LpOp", "list_key"},
			expected: "$1\r\na\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor := processor.NewProcessor()

			// Run setup commands
			for _, setupCmd := range tt.setup {
				processor.ProcessCommand(setupCmd)
			}

			// Run the actual LPOP test
			result := processor.ProcessCommand(tt.input)
			if result != tt.expected {
				t.Errorf("ProcessCommand(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func BenchmarkLPopCommand(b *testing.B) {
	testCases := []struct {
		name  string
		setup []string
		input []string
	}{
		{
			name:  "LPOP from small list",
			setup: []string{"RPUSH", "small_list", "a", "b", "c"},
			input: []string{"LPOP", "small_list"},
		},
		{
			name:  "LPOP from medium list",
			setup: []string{"RPUSH", "medium_list", "a", "b", "c", "d", "e", "f", "g"},
			input: []string{"LPOP", "medium_list"},
		},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			processor := processor.NewProcessor()

			// Setup the list for benchmarking
			processor.ProcessCommand(tc.setup)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				processor.ProcessCommand(tc.input)
			}
		})
	}
}
