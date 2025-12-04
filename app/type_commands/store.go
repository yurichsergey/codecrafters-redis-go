package type_commands

import (
	"github.com/codecrafters-io/redis-starter-go/app/list"
	"github.com/codecrafters-io/redis-starter-go/app/stream"
	"github.com/codecrafters-io/redis-starter-go/app/string_commands"
)

type Store struct {
	// StringStore holds reference to string storage for type checking
	StringStore *string_commands.Store
	// ListStore holds reference to list storage for type checking
	ListStore *list.Store
	// StreamStore holds reference to stream storage for type checking
	StreamStore *stream.Store
}

// NewStore creates a new Store instance with references to string, list, and stream stores.
func NewStore(stringStore *string_commands.Store, listStore *list.Store, streamStore *stream.Store) *Store {
	return &Store{
		StringStore: stringStore,
		ListStore:   listStore,
		StreamStore: streamStore,
	}
}
