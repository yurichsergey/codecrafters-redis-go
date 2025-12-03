package list

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

// LPush inserts one or more elements at the head of a list.
// Example: LPUSH mylist "world"
func (s *Store) LPush(row []string) string {
	// Check if there are enough arguments
	if len(row) < 3 {
		return resp.MakeError("ERR wrong number of arguments for 'lpush' command")
	}

	key := row[1]
	elements := row[2:]

	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Prepend elements in reverse order to maintain the correct list order
	for i, j := 0, len(elements)-1; i < j; i, j = i+1, j-1 {
		elements[i], elements[j] = elements[j], elements[i]
	}
	s.storage[key] = append(elements, s.storage[key]...)

	return resp.MakeInteger(len(s.storage[key]))
}
