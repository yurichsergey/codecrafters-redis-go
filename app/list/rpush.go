package list

import (
	"container/list"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

// RPush appends one or more elements to the end of a list.
// Example: RPUSH mylist "hello" "world"
func (s *Store) RPush(row []string) string {
	// Check if there are enough arguments
	if len(row) < 3 {
		return resp.MakeError("ERR wrong number of arguments for 'rpush' command")
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
	for _, element := range elements {
		l.PushBack(element)
	}

	// Calculate the new length of the list
	newLength := l.Len()

	if clients, exists := s.blockingClients[key]; exists {
		// Loop while we have both waiting clients and elements in the list
		for len(clients) > 0 && l.Len() > 0 {
			// Wake up the first (longest waiting) blocking client
			client := clients[0]

			// Get and remove the first element
			front := l.Front()
			val := front.Value.(string)
			l.Remove(front)

			client.Waiting <- BlockingResult{Key: key, Value: val}

			// Remove the first client
			clients = clients[1:]
		}

		// Update the blocking clients list
		s.blockingClients[key] = clients

		// Clean up empty list of blocking clients
		if len(s.blockingClients[key]) == 0 {
			delete(s.blockingClients, key)
		}
	}

	// Clean up empty list if all elements were consumed by blocking clients
	if l.Len() == 0 {
		delete(s.storage, key)
	}

	return resp.MakeInteger(newLength)
}
