package string_commands

import (
	"testing"
)

func TestSetCommand(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "SET with key and value",
			input:    []string{"SET", "foo", "bar"},
			expected: "+OK\r\n",
		},
		{
			name:     "SET with single character key and value",
			input:    []string{"SET", "a", "b"},
			expected: "+OK\r\n",
		},
		{
			name:     "SET with numeric value",
			input:    []string{"SET", "count", "123"},
			expected: "+OK\r\n",
		},
		{
			name:     "SET with empty value",
			input:    []string{"SET", "empty", ""},
			expected: "+OK\r\n",
		},
		{
			name:     "SET with spaces in value",
			input:    []string{"SET", "message", "hello world"},
			expected: "+OK\r\n",
		},
		{
			name:     "SET with special characters in value",
			input:    []string{"SET", "special", "!@#$%^&*()"},
			expected: "+OK\r\n",
		},
		{
			name:     "SET without key",
			input:    []string{"SET"},
			expected: "-ERR wrong number of arguments for 'set' command\r\n",
		},
		{
			name:     "SET with key but no value",
			input:    []string{"SET", "key"},
			expected: "-ERR wrong number of arguments for 'set' command\r\n",
		},
		{
			name:     "SET overwrites existing key",
			input:    []string{"SET", "mykey", "newvalue"},
			expected: "+OK\r\n",
		},
	}

	store := NewStore()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := store.Set(tt.input)
			if result != tt.expected {
				t.Errorf("Set(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSetCommandWithExpiry(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "SET with EX (expiry in seconds)",
			input:    []string{"SET", "key1", "value1", "EX", "10"},
			expected: "+OK\r\n",
		},
		{
			name:     "SET with PX (expiry in milliseconds)",
			input:    []string{"SET", "key2", "value2", "PX", "5000"},
			expected: "+OK\r\n",
		},
		{
			name:     "SET with EX zero seconds",
			input:    []string{"SET", "key3", "value3", "EX", "0"},
			expected: "+OK\r\n",
		},
		{
			name:     "SET with PX zero milliseconds",
			input:    []string{"SET", "key4", "value4", "PX", "0"},
			expected: "+OK\r\n",
		},
		{
			name:     "SET with invalid expiry type",
			input:    []string{"SET", "key5", "value5", "INVALID", "100"},
			expected: "-ERR syntax error\r\n",
		},
		{
			name:     "SET with non-numeric EX value",
			input:    []string{"SET", "key6", "value6", "EX", "abc"},
			expected: "-ERR value is not an integer or out of range\r\n",
		},
		{
			name:     "SET with non-numeric PX value",
			input:    []string{"SET", "key7", "value7", "PX", "xyz"},
			expected: "-ERR value is not an integer or out of range\r\n",
		},
		{
			name:     "SET with negative EX value",
			input:    []string{"SET", "key8", "value8", "EX", "-1"},
			expected: "+OK\r\n",
		},
		{
			name:     "SET with large EX value",
			input:    []string{"SET", "key9", "value9", "EX", "86400"},
			expected: "+OK\r\n",
		},
		{
			name:     "SET with large PX value",
			input:    []string{"SET", "key10", "value10", "PX", "86400000"},
			expected: "+OK\r\n",
		},
	}

	store := NewStore()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := store.Set(tt.input)
			if result != tt.expected {
				t.Errorf("Set(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
