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
	l, exists := s.storage[key]
	if !exists || l.Len() == 0 {
		return resp.MakeNullBulkString()
	}

	// Determine how many elements to actually remove
	numToRemove := count
	if numToRemove > l.Len() {
		numToRemove = l.Len()
	}

	// If count is 1 (no count argument provided), return single bulk string
	if len(row) == 2 {
		front := l.Front()
		val := front.Value.(string)
		l.Remove(front)

		// Clean up empty list
		if l.Len() == 0 {
			delete(s.storage, key)
		}

		return resp.MakeBulkString(val)
	}

	// Remove elements from the front
	removed := make([]string, 0, numToRemove)
	for i := 0; i < numToRemove; i++ {
		front := l.Front()
		val := front.Value.(string)
		removed = append(removed, val)
		l.Remove(front)
	}

	// Clean up empty list
	if l.Len() == 0 {
		delete(s.storage, key)
	}

	// Return as RESP array
	return resp.MakeArray(removed)
}
