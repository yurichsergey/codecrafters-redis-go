package stream

import (
	"sort"
)

// RadixTree is a compressed trie data structure specialized for stream entries.
type RadixTree struct {
	root *Node
	size int
}

// Node represents a node in the Radix Tree.
type Node struct {
	// prefix is the common prefix stored in this node
	prefix string
	// children maps the next character (byte) to the child node
	children []*Node
	// edges stores the first byte of the child's prefix for corresponding index in children
	// This keeps children sorted by edge byte.
	edges []byte
	// value holds the stream entry if this node represents a complete key
	value *Entry
}

// NewRadixTree creates a new empty Radix Tree.
func NewRadixTree() *RadixTree {
	return &RadixTree{
		root: &Node{},
	}
}

// Insert adds a stream entry to the tree using the given key (binary string).
func (t *RadixTree) Insert(key string, value *Entry) {
	node := t.root
	search := key

	for {
		// Find edge
		idx := -1
		if len(search) > 0 {
			firstChar := search[0]
			// Binary search for the edge
			n := len(node.edges)
			i := sort.Search(n, func(i int) bool { return node.edges[i] >= firstChar })
			if i < n && node.edges[i] == firstChar {
				idx = i
			}
		}

		// No matching edge found, insert remaining as new child
		if idx == -1 {
			if len(search) > 0 {
				newNode := &Node{
					prefix: search,
					value:  value,
				}
				t.addChild(node, newNode)
				t.size++
			} else {
				// Key ends at this node, update value
				// If value was nil, we are adding new key (though we might be just updating)
				// But strictly speaking for streams, we append unique IDs usually.
				// If key exists, just update.
				if node.value == nil {
					t.size++
				}
				node.value = value
			}
			return
		}

		// Edge found, traverse down
		child := node.children[idx]
		commonLen := commonPrefixLen(search, child.prefix)

		// Split child if common prefix is shorter than child's prefix
		if commonLen < len(child.prefix) {
			splitNode := &Node{
				prefix:   child.prefix[:commonLen],
				children: []*Node{child},
				edges:    []byte{child.prefix[commonLen]},
			}
			child.prefix = child.prefix[commonLen:]
			node.children[idx] = splitNode
			child = splitNode // Continue with the split node
		}

		search = search[commonLen:]

		// If search is exhausted, this is the node
		if len(search) == 0 {
			if child.value == nil {
				t.size++
			}
			child.value = value
			return
		}

		node = child
	}
}

// addChild adds a child node, keeping edges sorted.
func (t *RadixTree) addChild(parent *Node, child *Node) {
	edge := child.prefix[0]
	n := len(parent.edges)
	i := sort.Search(n, func(i int) bool { return parent.edges[i] >= edge })

	parent.edges = append(parent.edges, 0)
	copy(parent.edges[i+1:], parent.edges[i:])
	parent.edges[i] = edge

	parent.children = append(parent.children, nil)
	copy(parent.children[i+1:], parent.children[i:])
	parent.children[i] = child
}

// commonPrefixLen returns the length of the common prefix of a and b.
func commonPrefixLen(a, b string) int {
	i := 0
	for i < len(a) && i < len(b) && a[i] == b[i] {
		i++
	}
	return i
}

// Len returns the number of elements in the tree.
func (t *RadixTree) Len() int {
	return t.size
}

// Last returns the entry with the largest key (lexicographically).
func (t *RadixTree) Last() *Entry {
	node := t.root
	// We need to find the rightmost path
	// Since edges are sorted, the last child is always the "largest" path.

	// BUT, a node itself might have a value.
	// We need to traverse down taking the last child every time
	// AND checking if the node itself has a value.
	// Since keys in streams (time-seq) are strictly increasing and we store full 16-byte keys,
	// longer keys with same prefix are "larger" only if the suffix bytes are larger,
	// but here all keys are fixed 16 bytes.
	// So deeper nodes (children) are always strictly larger than parent if parent had a value
	// (which wouldn't happen for fixed length keys unless one is prefix of another,
	// but 16 byte keys cannot be prefix of another 16 byte key).

	// WAIT. Fixed length keys mean all values are at leaf nodes (or nodes at depth 16 bytes worth of prefix).
	// So just following the last child is sufficient.

	if len(node.children) == 0 && node.value == nil {
		return nil
	}

	for len(node.children) > 0 {
		node = node.children[len(node.children)-1]
	}

	return node.value
}

// First returns the entry with the smallest key (lexicographically).
func (t *RadixTree) First() *Entry {
	node := t.root
	if len(node.children) == 0 && node.value == nil {
		return nil
	}

	// Follow the first child
	// If a node has a value, it is a prefix of any children (if any).
	// If keys are fixed length (16 bytes), values are only at leaves (depth 16).
	// But in general Radix Tree, node value is "smaller" than children values?
	// Our keys are fixed length. So we just follow 0-th child until leaf.

	for len(node.children) > 0 {
		node = node.children[0]
	}

	return node.value
}

// Range returns all entries with keys in the range [start, end].
// start and end must be 16-byte strings.
func (t *RadixTree) Range(start, end string) []*Entry {
	var results []*Entry
	t.walk(t.root, "", start, end, &results)
	return results
}

func (t *RadixTree) walk(node *Node, path string, start, end string, results *[]*Entry) {
	currentPath := path + node.prefix

	// Check if current node has a value and is within range
	if node.value != nil {
		if currentPath >= start && currentPath <= end {
			*results = append(*results, node.value)
		}
	}

	// Optimization:
	// Since edges are sorted, children are sorted lexicographically.
	// We can potentially skip some children if we know they are out of range.
	// However, because nodes can condense prefixes, doing a precise check is tricky without full reconstruction.
	// A simple but effective optimization:
	// If currentPath is already > end, we might stop if we know we can't go back?
	// But currentPath is just a prefix. "prefix" > "end" is possible if prefix is "longer" or just lexicographically larger.
	// If currentPath is a prefix of 'end', we must continue.
	// If currentPath > end (lexicographically) and currentPath is NOT a prefix of end?
	//  e.g. currentPath="b", end="a". "b" > "a". All children of "b" will start with "b...", so they are > "a".
	// So if currentPath > end, we can prune THIS subtree.
	// Note: Strings in Go are compared byte by byte.

	if len(currentPath) > 0 {
		// If currentPath > end, then all children (which extend currentPath) are also > end.
		// Be careful: "10" < "2". But keys are fixed length binary strings (16 bytes),
		// so comparison is consistent.
		if currentPath > end {
			return
		}
	}

	// Also if we can determine that a child will be > end or < start?
	// Let's stick to the current path check for pruning > end.

	for _, child := range node.children {
		t.walk(child, currentPath, start, end, results)
	}
}
