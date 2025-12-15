package stream

import (
	"strings"
	"testing"
)

func TestXRange_Basic(t *testing.T) {
	store := NewStore()
	store.XAdd([]string{"XADD", "mystream", "100-1", "f1", "v1"})
	store.XAdd([]string{"XADD", "mystream", "100-2", "f2", "v2"})
	store.XAdd([]string{"XADD", "mystream", "100-3", "f3", "v3"})

	// Query range 100-1 to 100-2
	res := store.XRange([]string{"XRANGE", "mystream", "100-1", "100-2"})

	// Expect 2 entries
	// *2\r\n
	//   *2\r\n$5\r\n100-1\r\n*2\r\n$2\r\nf1\r\n$2\r\nv1\r\n
	//   *2\r\n$5\r\n100-2\r\n*2\r\n$2\r\nf2\r\n$2\r\nv2\r\n

	if !strings.Contains(res, "100-1") || !strings.Contains(res, "100-2") {
		t.Errorf("Expected 100-1 and 100-2, got %q", res)
	}
	if strings.Contains(res, "100-3") {
		t.Errorf("Did not expect 100-3, got %q", res)
	}

	// Count entries implicitly by checking start *2
	if !strings.HasPrefix(res, "*2\r\n") {
		t.Errorf("Expected array length 2, got %q", res)
	}
}

func TestXRange_PartialIDs(t *testing.T) {
	store := NewStore()
	store.XAdd([]string{"XADD", "s", "100-1", "a", "b"})
	store.XAdd([]string{"XADD", "s", "100-2", "c", "d"})
	store.XAdd([]string{"XADD", "s", "101-1", "e", "f"})

	// Range "100" "100" -> implies 100-0 to 100-MAX
	res := store.XRange([]string{"XRANGE", "s", "100", "100"})

	if !strings.Contains(res, "100-1") || !strings.Contains(res, "100-2") {
		t.Errorf("Expected 100-1 and 100-2, got %q", res)
	}
	if strings.Contains(res, "101-1") {
		t.Errorf("Did not expect 101-1, got %q", res)
	}
}

func TestXRange_MinMax(t *testing.T) {
	store := NewStore()
	store.XAdd([]string{"XADD", "s", "0-1", "a", "b"})
	store.XAdd([]string{"XADD", "s", "10-0", "c", "d"})
	store.XAdd([]string{"XADD", "s", "100-0", "e", "f"})

	// Range - +
	res := store.XRange([]string{"XRANGE", "s", "-", "+"})

	if !strings.HasPrefix(res, "*3\r\n") {
		t.Errorf("Expected 3 entries, got %q", res)
	}
}

func TestXRange_Empty(t *testing.T) {
	store := NewStore()
	res := store.XRange([]string{"XRANGE", "noonestream", "-", "+"})
	// Expect empty array
	if res != "*0\r\n" {
		t.Errorf("Expected empty array *0\\r\\n, got %q", res)
	}
}
