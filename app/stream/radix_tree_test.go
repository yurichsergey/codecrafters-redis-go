package stream

import (
	"testing"
)

func TestRadixTree_InsertAndLen(t *testing.T) {
	tree := NewRadixTree()

	if tree.Len() != 0 {
		t.Errorf("Expected empty tree to have len 0, got %d", tree.Len())
	}

	entry1 := &Entry{ID: "1-0"}
	key1, _ := IDToKey(entry1.ID)
	tree.Insert(key1, entry1)

	if tree.Len() != 1 {
		t.Errorf("Expected tree len 1, got %d", tree.Len())
	}

	entry2 := &Entry{ID: "1-1"}
	key2, _ := IDToKey(entry2.ID)
	tree.Insert(key2, entry2)

	if tree.Len() != 2 {
		t.Errorf("Expected tree len 2, got %d", tree.Len())
	}

	// Insert duplicate
	tree.Insert(key1, entry1)
	if tree.Len() != 2 {
		t.Errorf("Expected tree len 2 after duplicate insert, got %d", tree.Len())
	}
}

func TestRadixTree_Last(t *testing.T) {
	tree := NewRadixTree()

	if tree.Last() != nil {
		t.Error("Expected Last() to return nil for empty tree")
	}

	ids := []string{"1-0", "1-1", "0-1", "2-0", "10-0"}
	// Insert in random order? Or just mixed.
	// Expected last is 10-0

	// Convert checks
	for _, id := range ids {
		entry := &Entry{ID: id}
		key, _ := IDToKey(id)
		tree.Insert(key, entry)
	}

	last := tree.Last()
	if last == nil {
		t.Fatal("Expected Last() to return entry")
	}

	if last.ID != "10-0" {
		t.Errorf("Expected last ID to be 10-0, got %s", last.ID)
	}
}

func TestRadixTree_InternalStructure(t *testing.T) {
	// Glass-box testing to verify structure roughly
	tree := NewRadixTree()

	// Insert "A" (encoded)
	// Insert "AA" (encoded) -> Not possible with fixed length keys?
	// Oh right, our keys are ALL 16 bytes.
	// So prefixes will happen naturally.

	// Let's manually construct simpler keys for this test to verify splitting
	// We can't use IDToKey easily for custom simple strings, so let's bypass IDToKey for this specific test
	// But RadixTree uses string keys, so we can pass anything.

	// 1. Insert "apple"
	val1 := &Entry{ID: "apple"} // Dummy entry
	tree.Insert("apple", val1)

	if tree.root.children[0].prefix != "apple" {
		t.Errorf("Expected root child prefix 'apple', got '%s'", tree.root.children[0].prefix)
	}

	// 2. Insert "app" -> should split to "app" -> "le"
	val2 := &Entry{ID: "app"}
	tree.Insert("app", val2)

	// root -> "app" (val2) -> "le" (val1)
	if len(tree.root.children) != 1 {
		t.Fatalf("Expected 1 child of root, got %d", len(tree.root.children))
	}
	child := tree.root.children[0]
	if child.prefix != "app" {
		t.Errorf("Expected prefix 'app', got '%s'", child.prefix)
	}
	if child.value != val2 {
		t.Error("Expected 'app' node to have val2")
	}

	if len(child.children) != 1 {
		t.Fatalf("Expected 1 child of 'app' node, got %d", len(child.children))
	}
	grandChild := child.children[0]
	if grandChild.prefix != "le" {
		t.Errorf("Expected prefix 'le', got '%s'", grandChild.prefix)
	}
	if grandChild.value != val1 {
		t.Error("Expected 'le' node to have val1")
	}

	// 3. Insert "apricot" -> split "app" at "ap" -> "p" and "ricot"
	// root -> "ap" -> "p" (val2) -> "le" (val1)
	//              -> "ricot" (val3)
	val3 := &Entry{ID: "apricot"}
	tree.Insert("apricot", val3)

	if len(tree.root.children) != 1 {
		t.Fatalf("Expected 1 child of root, got %d", len(tree.root.children))
	}
	child = tree.root.children[0]
	if child.prefix != "ap" {
		t.Errorf("Expected prefix 'ap', got '%s'", child.prefix)
	}

	if len(child.children) != 2 {
		t.Fatalf("Expected 2 children of 'ap' node, got %d", len(child.children))
	}

	// Check edges order. 'p' vs 'r'. 'p' < 'r'.
	// So child 0 is 'p' (from app), child 1 is 'ricot'

	c0 := child.children[0]
	if c0.prefix != "p" {
		t.Errorf("Expected child 0 prefix 'p', got '%s'", c0.prefix)
	}
	if c0.value != val2 {
		t.Error("Expected 'p' node to have val2")
	}

	c1 := child.children[1]
	if c1.prefix != "ricot" {
		t.Errorf("Expected child 1 prefix 'ricot', got '%s'", c1.prefix)
	}
	if c1.value != val3 {
		t.Error("Expected 'ricot' node to have val3")
	}
}
