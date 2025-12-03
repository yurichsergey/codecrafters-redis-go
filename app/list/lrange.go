package list

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

// LRange returns the specified elements of the list stored at key.
// Example: LRANGE mylist 0 -1
func (s *Store) LRange(row []string) string {
	// Check if there are enough arguments
	if len(row) != 4 {
		return resp.MakeError("ERR wrong number of arguments for 'lrange' command")
	}

	key := row[1]
	startStr := row[2]
	stopStr := row[3]

	// Parse start index with negative index support
	start, err := strconv.Atoi(startStr)
	if err != nil {
		return resp.MakeError("ERR invalid start index")
	}

	// Parse stop index with negative index support
	stop, err := strconv.Atoi(stopStr)
	if err != nil {
		return resp.MakeError("ERR invalid stop index")
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Retrieve the list
	list, exists := s.storage[key]
	if !exists {
		// If list doesn't exist, return an empty array
		return resp.MakeEmptyArray()
	}

	// Calculate the length of the list
	listLength := len(list)

	// Handle negative indexes for start
	if start < 0 {
		start = listLength + start
		// If negative index is out of range, treat as 0
		if start < 0 {
			start = 0
		}
	}

	// Handle negative indexes for stop
	if stop < 0 {
		stop = listLength + stop
		// If negative index is out of range, treat as 0
		if stop < 0 {
			stop = 0
		}
	}

	// Adjust stop index if it's greater than or equal to list length
	if stop >= listLength {
		stop = listLength - 1
	}

	// Check if start index is out of bounds
	if start >= listLength {
		return resp.MakeEmptyArray()
	}

	// Check if start index is greater than stop index
	if start > stop {
		return resp.MakeEmptyArray()
	}

	// Extract the sublist
	subList := list[start : stop+1]

	// Construct the RESP array response
	return resp.MakeArray(subList)
}
