package list_test

import (
	"testing"

	"github.com/codecrafters-io/redis-starter-go/app/processor"
)

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
