package stream

import (
	"strings"
	"testing"
)

func TestXAdd_WildcardSequence(t *testing.T) {
	store := NewStore()

	// Scenario 1: Empty stream, time 0 -> 0-1
	id1 := store.XAdd([]string{"XADD", "stream1", "0-*", "f1", "v1"})
	if id1 != "$3\r\n0-1\r\n" {
		t.Errorf("Expected 0-1, got %s", id1)
	}

	// Scenario 2: Empty stream, time 1 -> 1-0
	id2 := store.XAdd([]string{"XADD", "stream2", "1-*", "f1", "v1"})
	if id2 != "$3\r\n1-0\r\n" {
		t.Errorf("Expected 1-0, got %s", id2)
	}

	// Scenario 3: Stream with 1-0, add 1-* -> 1-1
	id3 := store.XAdd([]string{"XADD", "stream2", "1-*", "f1", "v1"})
	if id3 != "$3\r\n1-1\r\n" {
		t.Errorf("Expected 1-1, got %s", id3)
	}

	// Scenario 4: Stream with 1-1, add 2-* -> 2-0
	id4 := store.XAdd([]string{"XADD", "stream2", "2-*", "f1", "v1"})
	if id4 != "$3\r\n2-0\r\n" {
		t.Errorf("Expected 2-0, got %s", id4)
	}

	// Scenario 5: Stream with 2-0, add 0-* -> Error (0-1 <= 2-0)
	errResp := store.XAdd([]string{"XADD", "stream2", "0-*", "f1", "v1"})
	if !strings.HasPrefix(errResp, "-ERR") {
		t.Errorf("Expected error for Time < LastTime, got %s", errResp)
	}
}
