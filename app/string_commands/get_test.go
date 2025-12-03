package string_commands

import (
	"testing"
	"time"
)

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
			store := NewStore()

			// Run setup commands
			for _, setupCmd := range tt.setup {
				if setupCmd[0] == "SET" {
					store.Set(setupCmd)
				}
			}

			// Run the actual test
			result := store.Get(tt.input)
			if result != tt.expected {
				t.Errorf("Get(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetCommandWithExpiry(t *testing.T) {
	t.Run("GET key before EX expiry", func(t *testing.T) {
		store := NewStore()

		// SET with 2 second expiry
		store.Set([]string{"SET", "tempkey", "tempvalue", "EX", "2"})

		// GET immediately should return the value
		result := store.Get([]string{"GET", "tempkey"})
		expected := "$9\r\ntempvalue\r\n"
		if result != expected {
			t.Errorf("GET before expiry failed: got %q, want %q", result, expected)
		}
	})

	t.Run("GET key after EX expiry", func(t *testing.T) {
		store := NewStore()

		// SET with the 1 millisecond expiry using PX
		store.Set([]string{"SET", "tempkey", "tempvalue", "PX", "1"})

		// Wait for expiry
		time.Sleep(10 * time.Millisecond)

		// GET after expiry should return null
		result := store.Get([]string{"GET", "tempkey"})
		expected := "$-1\r\n"
		if result != expected {
			t.Errorf("GET after expiry failed: got %q, want %q", result, expected)
		}
	})

	t.Run("GET key before PX expiry", func(t *testing.T) {
		store := NewStore()

		// SET with 1000ms (1 second) expiry
		store.Set([]string{"SET", "pxkey", "pxvalue", "PX", "1000"})

		// GET immediately should return the value
		result := store.Get([]string{"GET", "pxkey"})
		expected := "$7\r\npxvalue\r\n"
		if result != expected {
			t.Errorf("GET before PX expiry failed: got %q, want %q", result, expected)
		}
	})

	t.Run("GET key after PX expiry", func(t *testing.T) {
		store := NewStore()

		// SET with 50 millisecond expiry
		store.Set([]string{"SET", "pxkey", "pxvalue", "PX", "50"})

		// Wait for expiry
		time.Sleep(100 * time.Millisecond)

		// GET after expiry should return null
		result := store.Get([]string{"GET", "pxkey"})
		expected := "$-1\r\n"
		if result != expected {
			t.Errorf("GET after PX expiry failed: got %q, want %q", result, expected)
		}
	})

	t.Run("SET without expiry does not expire", func(t *testing.T) {
		store := NewStore()

		// SET without expiry
		store.Set([]string{"SET", "noexpiry", "persistent"})

		// Wait some time
		time.Sleep(50 * time.Millisecond)

		// GET should still return the value
		result := store.Get([]string{"GET", "noexpiry"})
		expected := "$10\r\npersistent\r\n"
		if result != expected {
			t.Errorf("GET non-expiring key failed: got %q, want %q", result, expected)
		}
	})

	t.Run("Overwrite key with new expiry", func(t *testing.T) {
		store := NewStore()

		// SET with short expiry
		store.Set([]string{"SET", "overwrite", "oldvalue", "PX", "50"})

		// Immediately overwrite with longer expiry
		store.Set([]string{"SET", "overwrite", "newvalue", "EX", "10"})

		// Wait past the first expiry time
		time.Sleep(100 * time.Millisecond)

		// Key should still exist with new value
		result := store.Get([]string{"GET", "overwrite"})
		expected := "$8\r\nnewvalue\r\n"
		if result != expected {
			t.Errorf("GET overwritten key failed: got %q, want %q", result, expected)
		}
	})
}
