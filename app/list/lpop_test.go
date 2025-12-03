package list_test

import (
	"testing"

	"github.com/codecrafters-io/redis-starter-go/app/processor"
)

func TestLPopSingleElement(t *testing.T) {
	p := processor.NewProcessor()

	// Setup: Create a list with multiple elements
	response := p.ProcessCommand([]string{"RPUSH", "list_key", "one", "two", "three"})
	expected := ":3\r\n"
	if response != expected {
		t.Errorf("RPUSH failed. Expected %q, got %q", expected, response)
	}

	// Test: LPOP without count (should return single bulk string)
	response = p.ProcessCommand([]string{"LPOP", "list_key"})
	expected = "$3\r\none\r\n"
	if response != expected {
		t.Errorf("LPOP single element failed. Expected %q, got %q", expected, response)
	}

	// Verify remaining elements
	response = p.ProcessCommand([]string{"LRANGE", "list_key", "0", "-1"})
	expected = "*2\r\n$3\r\ntwo\r\n$5\r\nthree\r\n"
	if response != expected {
		t.Errorf("LRANGE after LPOP failed. Expected %q, got %q", expected, response)
	}
}

func TestLPopMultipleElements(t *testing.T) {
	p := processor.NewProcessor()

	// Setup: Create a list
	p.ProcessCommand([]string{"RPUSH", "list_key", "a", "b", "c", "d"})

	// Test: LPOP with count 2
	response := p.ProcessCommand([]string{"LPOP", "list_key", "2"})
	expected := "*2\r\n$1\r\na\r\n$1\r\nb\r\n"
	if response != expected {
		t.Errorf("LPOP with count 2 failed. Expected %q, got %q", expected, response)
	}

	// Verify remaining elements
	response = p.ProcessCommand([]string{"LRANGE", "list_key", "0", "-1"})
	expected = "*2\r\n$1\r\nc\r\n$1\r\nd\r\n"
	if response != expected {
		t.Errorf("LRANGE after LPOP failed. Expected %q, got %q", expected, response)
	}
}

func TestLPopCountGreaterThanLength(t *testing.T) {
	p := processor.NewProcessor()

	// Setup: Create a list with 3 elements
	p.ProcessCommand([]string{"RPUSH", "list_key", "one", "two", "three"})

	// Test: LPOP with count 10 (greater than list length)
	response := p.ProcessCommand([]string{"LPOP", "list_key", "10"})
	expected := "*3\r\n$3\r\none\r\n$3\r\ntwo\r\n$5\r\nthree\r\n"
	if response != expected {
		t.Errorf("LPOP with count > length failed. Expected %q, got %q", expected, response)
	}

	// Verify list is empty
	response = p.ProcessCommand([]string{"LRANGE", "list_key", "0", "-1"})
	expected = "*0\r\n"
	if response != expected {
		t.Errorf("List should be empty. Expected %q, got %q", expected, response)
	}
}

func TestLPopEmptyList(t *testing.T) {
	p := processor.NewProcessor()

	// Test: LPOP on non-existent list
	response := p.ProcessCommand([]string{"LPOP", "non_existent"})
	expected := "$-1\r\n"
	if response != expected {
		t.Errorf("LPOP on empty list failed. Expected %q, got %q", expected, response)
	}

	// Test: LPOP with count on non-existent list
	response = p.ProcessCommand([]string{"LPOP", "non_existent", "2"})
	if response != expected {
		t.Errorf("LPOP with count on empty list failed. Expected %q, got %q", expected, response)
	}
}

func TestLPopCountOne(t *testing.T) {
	p := processor.NewProcessor()

	// Setup
	p.ProcessCommand([]string{"RPUSH", "list_key", "alpha", "beta", "gamma"})

	// Test: LPOP with explicit count 1 (should return array with 1 element)
	response := p.ProcessCommand([]string{"LPOP", "list_key", "1"})
	expected := "*1\r\n$5\r\nalpha\r\n"
	if response != expected {
		t.Errorf("LPOP with count 1 failed. Expected %q, got %q", expected, response)
	}

	// Verify remaining elements
	response = p.ProcessCommand([]string{"LLEN", "list_key"})
	expected = ":2\r\n"
	if response != expected {
		t.Errorf("LLEN after LPOP failed. Expected %q, got %q", expected, response)
	}
}

func TestLPopAllElements(t *testing.T) {
	p := processor.NewProcessor()

	// Setup
	p.ProcessCommand([]string{"RPUSH", "list_key", "one", "two", "three", "four", "five"})

	// Test: LPOP exactly all elements
	response := p.ProcessCommand([]string{"LPOP", "list_key", "5"})
	expected := "*5\r\n$3\r\none\r\n$3\r\ntwo\r\n$5\r\nthree\r\n$4\r\nfour\r\n$4\r\nfive\r\n"
	if response != expected {
		t.Errorf("LPOP all elements failed. Expected %q, got %q", expected, response)
	}

	// Verify list is deleted
	response = p.ProcessCommand([]string{"LLEN", "list_key"})
	expected = ":0\r\n"
	if response != expected {
		t.Errorf("List should be empty. Expected %q, got %q", expected, response)
	}
}

func TestLPopInvalidCount(t *testing.T) {
	p := processor.NewProcessor()

	// Setup
	p.ProcessCommand([]string{"RPUSH", "list_key", "a", "b", "c"})

	// Test: LPOP with negative count
	response := p.ProcessCommand([]string{"LPOP", "list_key", "-1"})
	if !contains(response, "ERR") {
		t.Errorf("LPOP with negative count should return error. Got %q", response)
	}

	// Test: LPOP with non-integer count
	response = p.ProcessCommand([]string{"LPOP", "list_key", "abc"})
	if !contains(response, "ERR") {
		t.Errorf("LPOP with non-integer count should return error. Got %q", response)
	}
}

func TestLPopSequentialOperations(t *testing.T) {
	p := processor.NewProcessor()

	// Setup
	p.ProcessCommand([]string{"RPUSH", "list_key", "1", "2", "3", "4", "5", "6"})

	// First LPOP
	response := p.ProcessCommand([]string{"LPOP", "list_key", "2"})
	expected := "*2\r\n$1\r\n1\r\n$1\r\n2\r\n"
	if response != expected {
		t.Errorf("First LPOP failed. Expected %q, got %q", expected, response)
	}

	// Second LPOP
	response = p.ProcessCommand([]string{"LPOP", "list_key", "3"})
	expected = "*3\r\n$1\r\n3\r\n$1\r\n4\r\n$1\r\n5\r\n"
	if response != expected {
		t.Errorf("Second LPOP failed. Expected %q, got %q", expected, response)
	}

	// Third LPOP (only 1 element left)
	response = p.ProcessCommand([]string{"LPOP", "list_key", "2"})
	expected = "*1\r\n$1\r\n6\r\n"
	if response != expected {
		t.Errorf("Third LPOP failed. Expected %q, got %q", expected, response)
	}

	// Verify list is empty
	response = p.ProcessCommand([]string{"LLEN", "list_key"})
	expected = ":0\r\n"
	if response != expected {
		t.Errorf("List should be empty. Expected %q, got %q", expected, response)
	}
}

func TestLPopZeroCount(t *testing.T) {
	p := processor.NewProcessor()

	// Setup
	p.ProcessCommand([]string{"RPUSH", "list_key", "a", "b", "c"})

	// Test: LPOP with count 0
	response := p.ProcessCommand([]string{"LPOP", "list_key", "0"})
	expected := "*0\r\n"
	if response != expected {
		t.Errorf("LPOP with count 0 failed. Expected %q, got %q", expected, response)
	}

	// Verify all elements remain
	response = p.ProcessCommand([]string{"LLEN", "list_key"})
	expected = ":3\r\n"
	if response != expected {
		t.Errorf("All elements should remain. Expected %q, got %q", expected, response)
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
