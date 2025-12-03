package list

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

// LLen returns the length of the list stored at key.
// Example: LLEN mylist
func (s *Store) LLen(row []string) string {
	// Check if there are enough arguments
	if len(row) != 2 {
		return resp.MakeError("ERR wrong number of arguments for 'llen' command")
	}

	key := row[1]

	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Get the list length
	list, exists := s.storage[key]
	if !exists {
		// If list doesn't exist, return 0
		return resp.MakeInteger(0)
	}

	// Return the length of the list as a RESP integer
	return resp.MakeInteger(len(list))
}
