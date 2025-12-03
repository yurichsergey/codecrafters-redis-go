package string_commands

import (
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

// Get gets the value of a key.
// Example: GET mykey
func (s *Store) Get(args []string) string {
	// GET command requires a key argument
	if len(args) < 2 {
		return resp.MakeError("ERR wrong number of arguments for 'get' command")
	}

	key := args[1]

	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if the key exists in storage
	item, exists := s.storage[key]
	if !exists {
		// Return null bulk string if the key doesn't exist
		return resp.MakeNullBulkString()
	}

	if item.Expiry != 0 && time.Now().UnixMilli() > item.Expiry {
		// delete(s.storage, key) // Lazy deletion could be implemented here
		return resp.MakeNullBulkString()
	}

	// Return the value as a RESP bulk string
	// Format: $<length>\r\n<data>\r\n
	return resp.MakeBulkString(item.Value)
}
