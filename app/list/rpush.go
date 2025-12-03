package list

import (
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

	s.storage[key] = append(s.storage[key], elements...)

	// Calculate the new length of the list
	newLength := len(s.storage[key])

	if clients, exists := s.blockingClients[key]; exists {
		// Loop while we have both waiting clients and elements in the list
		for len(clients) > 0 && len(s.storage[key]) > 0 {
			// Wake up the first (longest waiting) blocking client
			client := clients[0]
			client.Waiting <- BlockingResult{Key: key, Value: s.storage[key][0]}

			// Remove the first element and the first client
			s.storage[key] = s.storage[key][1:]
			clients = clients[1:]
		}

		// Update the blocking clients list
		s.blockingClients[key] = clients

		// Clean up empty list of blocking clients
		if len(s.blockingClients[key]) == 0 {
			delete(s.blockingClients, key)
		}
	}

	return resp.MakeInteger(newLength)
}
