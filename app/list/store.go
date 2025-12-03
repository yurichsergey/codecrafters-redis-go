package list

import (
	"sync"
)

type BlockingResult struct {
	// key is the list key that the client was waiting for
	Key string
	// value is the element popped from the list
	Value string
}

type BlockingClient struct {
	// Waiting is a channel that receives the result when an element is available
	Waiting chan BlockingResult
}

type Store struct {
	// storage holds the key-value pairs for list commands
	storage map[string][]string
	// blockingClients holds the list of clients waiting for elements on specific keys
	blockingClients map[string][]*BlockingClient
	// mutex protects access to the storage and blockingClients map
	mutex sync.Mutex
}

// NewStore creates a new Store instance with initialized storage and blocking clients.
func NewStore() *Store {
	return &Store{
		storage:         make(map[string][]string),
		blockingClients: make(map[string][]*BlockingClient),
	}
}
