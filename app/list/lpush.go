package list

import (
	"container/list"

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

	// Initialize list if it doesn't exist
	if _, exists := s.storage[key]; !exists {
		s.storage[key] = list.New()
	}

	l := s.storage[key]
	// Prepend elements. Redis LPUSH appends elements to the head.
	// LPUSH mylist A B C -> C is head, B second, A third.
	// So we push A, then B, then C to Front.
	// Wait, code says: elements are inserted one after the other to the head.
	// If I push A (Front: A), then B (Front: B, A), then C (Front: C, B, A).
	// So yes, iterating normally and PushFront works.
	for _, element := range elements {
		l.PushFront(element)
	}

	// Original LPush did NOT handle blocking clients. Assuming this is intended for now.
	// If we needed to handle them, we would do it here.

	return resp.MakeInteger(l.Len())
}
