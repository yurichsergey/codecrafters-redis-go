package list_test

import (
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
