package string_commands

import (
	"sync"
)

type StorageItem struct {
	// value is the string value stored in the item
	Value string
	// expiry is the expiration time in milliseconds
	Expiry int64
}

type Store struct {
	// storage holds the key-value pairs for string commands
	storage map[string]*StorageItem
	// mutex protects access to the storage map
	mutex sync.Mutex
}

// NewStore creates a new Store instance with initialized storage.
func NewStore() *Store {
	return &Store{
		storage: make(map[string]*StorageItem),
	}
}
