package list

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

// LPop removes and returns the first elements of the list stored at key.
// Example: LPOP mylist
func (s *Store) LPop(row []string) string {
	if len(row) < 2 {
		return resp.MakeError("ERR wrong number of arguments for 'lpop' command")
	}

	key := row[1]

	// Default count is 1
	count := 1

	// If count argument is provided
	if len(row) >= 3 {
		var err error
		count, err = strconv.Atoi(row[2])
		if err != nil || count < 0 {
			return resp.MakeError("ERR value is not an integer or out of range")
		}
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if list exists
	list, exists := s.storage[key]
	if !exists || len(list) == 0 {
		return resp.MakeNullBulkString()
	}

	// Determine how many elements to actually remove
	numToRemove := count
	if numToRemove > len(list) {
		numToRemove = len(list)
	}

	// If count is 1 (no count argument provided), return single bulk string
	if len(row) == 2 {
		removed := list[0]
		s.storage[key] = list[1:]

		// Clean up empty list
		if len(s.storage[key]) == 0 {
			delete(s.storage, key)
		}

		return resp.MakeBulkString(removed)
	}

	// Remove elements from the front
	removed := list[:numToRemove]
	s.storage[key] = list[numToRemove:]

	// Clean up empty list
	if len(s.storage[key]) == 0 {
		delete(s.storage, key)
	}

	// Return as RESP array
	return resp.MakeArray(removed)
}
