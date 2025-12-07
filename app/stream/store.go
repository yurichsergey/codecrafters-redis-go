package stream

import (
	"sync"
)

type Entry struct {
	// ID is the unique identifier for the entry (e.g., "1526919030474-0")
	ID string
	// Fields holds the key-value pairs for this entry
	Fields map[string]string
}

type Stream struct {
	// tree holds all entries in the stream sorted by ID
	tree *RadixTree
}

type Store struct {
	// storage maps stream keys to their corresponding streams
	storage map[string]*Stream
	// mutex protects concurrent access to the storage map
	mutex sync.Mutex
}

// NewStore creates a new Store instance with initialized storage.
func NewStore() *Store {
	return &Store{
		storage: make(map[string]*Stream),
	}
}
