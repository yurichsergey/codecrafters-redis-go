package list

import (
	"strconv"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

// BLPop removes and returns the first element of the list stored at key,
// blocking if the list is empty.
// Example: BLPOP mylist 0
func (s *Store) BLPop(row []string) string {
	// Check if there are enough arguments
	if len(row) < 3 {
		return resp.MakeError("ERR wrong number of arguments for 'blpop' command")
	}

	// Parse timeout as float
	timeoutStr := row[len(row)-1]
	timeoutSeconds, err := strconv.ParseFloat(timeoutStr, 64)
	if err != nil || timeoutSeconds < 0 {
		return resp.MakeError("ERR timeout is not a float or out of range")
	}

	// Get the keys
	keys := row[1 : len(row)-1]

	s.mutex.Lock()
	// Check each list for an element
	for _, key := range keys {
		list, exists := s.storage[key]
		if exists && len(list) > 0 {
			// Pop the first element
			element := list[0]
			s.storage[key] = list[1:]

			// Clean up empty list
			if len(s.storage[key]) == 0 {
				delete(s.storage, key)
			}

			s.mutex.Unlock()
			// Return the key and element as a RESP array
			return resp.MakeArray([]string{key, element})
		}
	}

	// If no elements are available, create a blocking client
	blockingClient := &BlockingClient{
		Waiting: make(chan BlockingResult, 1),
	}

	for _, key := range keys {
		s.blockingClients[key] = append(s.blockingClients[key], blockingClient)
	}
	s.mutex.Unlock()

	// Define a cleanup function to remove the client from all keys
	cleanup := func() {
		s.mutex.Lock()
		defer s.mutex.Unlock()
		for _, key := range keys {
			clients := s.blockingClients[key]
			for i, client := range clients {
				if client == blockingClient {
					// Remove the client
					s.blockingClients[key] = append(clients[:i], clients[i+1:]...)
					break
				}
			}
			// Clean up empty client lists
			if len(s.blockingClients[key]) == 0 {
				delete(s.blockingClients, key)
			}
		}
	}
	defer cleanup()

	// Block until an element is available or timeout expires
	var result BlockingResult
	if timeoutSeconds == 0 {
		// Indefinite blocking
		result = <-blockingClient.Waiting
	} else {
		// Blocking with timeout
		timer := time.NewTimer(time.Duration(timeoutSeconds * float64(time.Second)))
		defer timer.Stop()

		select {
		case result = <-blockingClient.Waiting:
			// Element received
		case <-timer.C:
			// Timeout expired
			return resp.MakeNullArray()
		}
	}

	// Return the key and element
	return resp.MakeArray([]string{result.Key, result.Value})
}
