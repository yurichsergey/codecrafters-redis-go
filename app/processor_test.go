package main

import (
	"testing"
	"time"
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
			input:    []string{"UNKNOWN", "key"},
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

func TestGetCommand(t *testing.T) {
	tests := []struct {
		name     string
		setup    [][]string // Commands to run before the test
		input    []string
		expected string
	}{
		{
			name:     "GET non-existent key",
			setup:    nil,
			input:    []string{"GET", "nonexistent"},
			expected: "$-1\r\n",
		},
		{
			name:     "GET existing key",
			setup:    [][]string{{"SET", "foo", "bar"}},
			input:    []string{"GET", "foo"},
			expected: "$3\r\nbar\r\n",
		},
		{
			name:     "GET single character value",
			setup:    [][]string{{"SET", "a", "x"}},
			input:    []string{"GET", "a"},
			expected: "$1\r\nx\r\n",
		},
		{
			name:     "GET numeric value",
			setup:    [][]string{{"SET", "count", "12345"}},
			input:    []string{"GET", "count"},
			expected: "$5\r\n12345\r\n",
		},
		{
			name:     "GET empty value",
			setup:    [][]string{{"SET", "empty", ""}},
			input:    []string{"GET", "empty"},
			expected: "$0\r\n\r\n",
		},
		{
			name:     "GET value with spaces",
			setup:    [][]string{{"SET", "message", "hello world"}},
			input:    []string{"GET", "message"},
			expected: "$11\r\nhello world\r\n",
		},
		{
			name:     "GET value with special characters",
			setup:    [][]string{{"SET", "special", "!@#$%"}},
			input:    []string{"GET", "special"},
			expected: "$5\r\n!@#$%\r\n",
		},
		{
			name:     "GET without key argument",
			setup:    nil,
			input:    []string{"GET"},
			expected: "-ERR wrong number of arguments for 'get' command\r\n",
		},
		{
			name:     "GET after overwriting value",
			setup:    [][]string{{"SET", "mykey", "oldvalue"}, {"SET", "mykey", "newvalue"}},
			input:    []string{"GET", "mykey"},
			expected: "$8\r\nnewvalue\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor := NewProcessor()

			// Run setup commands
			for _, setupCmd := range tt.setup {
				processor.ProcessCommand(setupCmd)
			}

			// Run the actual test
			result := processor.ProcessCommand(tt.input)
			if result != tt.expected {
				t.Errorf("ProcessCommand(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSetGetIntegration(t *testing.T) {
	processor := NewProcessor()

	t.Run("SET and GET same key", func(t *testing.T) {
		// Set a value
		setResult := processor.ProcessCommand([]string{"SET", "testkey", "testvalue"})
		if setResult != "+OK\r\n" {
			t.Errorf("SET failed: got %q, want %q", setResult, "+OK\r\n")
		}

		// Get the value back
		getResult := processor.ProcessCommand([]string{"GET", "testkey"})
		expected := "$9\r\ntestvalue\r\n"
		if getResult != expected {
			t.Errorf("GET failed: got %q, want %q", getResult, expected)
		}
	})

	t.Run("SET multiple keys and GET them", func(t *testing.T) {
		keys := map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		}

		// Set all keys
		for key, value := range keys {
			result := processor.ProcessCommand([]string{"SET", key, value})
			if result != "+OK\r\n" {
				t.Errorf("SET %s failed: got %q", key, result)
			}
		}

		// Get all keys back
		for key, value := range keys {
			result := processor.ProcessCommand([]string{"GET", key})
			expected := "$" + string(rune(len(value)+48)) + "\r\n" + value + "\r\n"
			if result != expected {
				t.Errorf("GET %s failed: got %q, want %q", key, result, expected)
			}
		}
	})

	t.Run("Overwrite key and verify new value", func(t *testing.T) {
		// Set initial value
		processor.ProcessCommand([]string{"SET", "updatekey", "oldvalue"})

		// Overwrite with new value
		setResult := processor.ProcessCommand([]string{"SET", "updatekey", "newvalue"})
		if setResult != "+OK\r\n" {
			t.Errorf("SET failed: got %q", setResult)
		}

		// Verify new value
		getResult := processor.ProcessCommand([]string{"GET", "updatekey"})
		expected := "$8\r\nnewvalue\r\n"
		if getResult != expected {
			t.Errorf("GET failed: got %q, want %q", getResult, expected)
		}
	})
}

func BenchmarkDefineResponse(b *testing.B) {
	testCases := []struct {
		name  string
		input []string
	}{
		{"PING", []string{"PING"}},
		{"ECHO single", []string{"ECHO", "hello"}},
		{"ECHO multiple", []string{"ECHO", "hello", "world", "test"}},
		{"SET", []string{"SET", "key", "value"}},
		{"GET existing", []string{"GET", "key"}},
	}

	processor := NewProcessor()
	// Setup for GET benchmark
	processor.ProcessCommand([]string{"SET", "key", "value"})

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				processor.ProcessCommand(tc.input)
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

func TestGetCommandWithExpiry(t *testing.T) {
	t.Run("GET key before EX expiry", func(t *testing.T) {
		processor := NewProcessor()

		// SET with 2 second expiry
		processor.ProcessCommand([]string{"SET", "tempkey", "tempvalue", "EX", "2"})

		// GET immediately should return the value
		result := processor.ProcessCommand([]string{"GET", "tempkey"})
		expected := "$9\r\ntempvalue\r\n"
		if result != expected {
			t.Errorf("GET before expiry failed: got %q, want %q", result, expected)
		}
	})

	t.Run("GET key after EX expiry", func(t *testing.T) {
		processor := NewProcessor()

		// SET with 1 millisecond expiry using PX
		processor.ProcessCommand([]string{"SET", "tempkey", "tempvalue", "PX", "1"})

		// Wait for expiry
		time.Sleep(10 * time.Millisecond)

		// GET after expiry should return null
		result := processor.ProcessCommand([]string{"GET", "tempkey"})
		expected := "$-1\r\n"
		if result != expected {
			t.Errorf("GET after expiry failed: got %q, want %q", result, expected)
		}
	})

	t.Run("GET key before PX expiry", func(t *testing.T) {
		processor := NewProcessor()

		// SET with 1000ms (1 second) expiry
		processor.ProcessCommand([]string{"SET", "pxkey", "pxvalue", "PX", "1000"})

		// GET immediately should return the value
		result := processor.ProcessCommand([]string{"GET", "pxkey"})
		expected := "$7\r\npxvalue\r\n"
		if result != expected {
			t.Errorf("GET before PX expiry failed: got %q, want %q", result, expected)
		}
	})

	t.Run("GET key after PX expiry", func(t *testing.T) {
		processor := NewProcessor()

		// SET with 50 millisecond expiry
		processor.ProcessCommand([]string{"SET", "pxkey", "pxvalue", "PX", "50"})

		// Wait for expiry
		time.Sleep(100 * time.Millisecond)

		// GET after expiry should return null
		result := processor.ProcessCommand([]string{"GET", "pxkey"})
		expected := "$-1\r\n"
		if result != expected {
			t.Errorf("GET after PX expiry failed: got %q, want %q", result, expected)
		}
	})

	t.Run("SET without expiry does not expire", func(t *testing.T) {
		processor := NewProcessor()

		// SET without expiry
		processor.ProcessCommand([]string{"SET", "noexpiry", "persistent"})

		// Wait some time
		time.Sleep(50 * time.Millisecond)

		// GET should still return the value
		result := processor.ProcessCommand([]string{"GET", "noexpiry"})
		expected := "$10\r\npersistent\r\n"
		if result != expected {
			t.Errorf("GET non-expiring key failed: got %q, want %q", result, expected)
		}
	})

	t.Run("Overwrite key with new expiry", func(t *testing.T) {
		processor := NewProcessor()

		// SET with short expiry
		processor.ProcessCommand([]string{"SET", "overwrite", "oldvalue", "PX", "50"})

		// Immediately overwrite with longer expiry
		processor.ProcessCommand([]string{"SET", "overwrite", "newvalue", "EX", "10"})

		// Wait past the first expiry time
		time.Sleep(100 * time.Millisecond)

		// Key should still exist with new value
		result := processor.ProcessCommand([]string{"GET", "overwrite"})
		expected := "$8\r\nnewvalue\r\n"
		if result != expected {
			t.Errorf("GET overwritten key failed: got %q, want %q", result, expected)
		}
	})
}
