package stream

import (
	"fmt"
	"strings"
	"testing"
)

func parseBulkString(s string) (string, error) {
	lines := strings.Split(s, "\r\n")
	if len(lines) < 2 {
		return "", fmt.Errorf("invalid bulk string format: %s", s)
	}
	return lines[1], nil
}

func TestXAdd_FullWildcard(t *testing.T) {
	store := NewStore()

	// 1. Add first entry with *
	resp1 := store.XAdd([]string{"XADD", "stream_key", "*", "foo", "bar"})
	if strings.HasPrefix(resp1, "-ERR") {
		t.Fatalf("XADD * returned error: %s", resp1)
	}

	id1, err := parseBulkString(resp1)
	if err != nil {
		t.Fatalf("Failed to parse RESP: %v", err)
	}

	// Parse ID1
	msTime1, seq1, err := ParseID(id1)
	if err != nil {
		t.Fatalf("Failed to parse first ID %s: %v", id1, err)
	}
	if msTime1 <= 0 {
		t.Errorf("Time part should be > 0, got %d", msTime1)
	}
	if seq1 < 0 {
		t.Errorf("Sequence part should be >= 0, got %d", seq1)
	}

	// 2. Add second entry with *
	resp2 := store.XAdd([]string{"XADD", "stream_key", "*", "baz", "qux"})
	if strings.HasPrefix(resp2, "-ERR") {
		t.Fatalf("XADD * returned error for second entry: %s", resp2)
	}
	id2, err := parseBulkString(resp2)
	if err != nil {
		t.Fatalf("Failed to parse RESP: %v", err)
	}

	msTime2, seq2, err := ParseID(id2)
	if err != nil {
		t.Fatalf("Failed to parse second ID %s: %v", id2, err)
	}

	// 3. Verify monotonicity
	// id2 > id1
	isGreater := false
	if msTime2 > msTime1 {
		isGreater = true
	} else if msTime2 == msTime1 {
		if seq2 > seq1 {
			isGreater = true
		}
	}
	if !isGreater {
		t.Errorf("Second ID (%s) should be strictly greater than First ID (%s)", id2, id1)
	}
}

func TestXAdd_FullWildcard_Collision(t *testing.T) {
	store := NewStore()
	var lastID string

	for i := 0; i < 100; i++ {
		respStr := store.XAdd([]string{"XADD", "collision_stream", "*", "k", "v"})
		if strings.HasPrefix(respStr, "-ERR") {
			t.Fatalf("XADD * iteration %d failed: %s", i, respStr)
		}

		id, err := parseBulkString(respStr)
		if err != nil {
			t.Fatalf("Failed to parse RESP: %v", err)
		}

		if lastID != "" {
			msLast, seqLast, _ := ParseID(lastID)
			msCurr, seqCurr, _ := ParseID(id)

			isGreater := (msCurr > msLast) || (msCurr == msLast && seqCurr > seqLast)
			if !isGreater {
				t.Fatalf("New ID %s not greater than Last ID %s", id, lastID)
			}
		}
		lastID = id
	}
}
