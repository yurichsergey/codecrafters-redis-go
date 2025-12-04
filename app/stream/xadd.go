package stream

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

// XAdd appends an entry to a stream and returns the entry ID.
// Example: XAdd(["XADD", "mystream", "0-1", "temperature", "36", "humidity", "95"])
func (s *Store) XAdd(args []string) string {
	// XADD requires at least: command, key, ID, and one field-value pair
	if len(args) < 5 {
		return resp.MakeError("ERR wrong number of arguments for 'xadd' command")
	}

	// Field-value pairs must come in pairs (even number of arguments after ID)
	if (len(args)-3)%2 != 0 {
		return resp.MakeError("ERR wrong number of arguments for 'xadd' command")
	}

	key := args[1]
	entryID := args[2]

	// Parse field-value pairs
	fields := make(map[string]string)
	for i := 3; i < len(args); i += 2 {
		fieldName := args[i]
		fieldValue := args[i+1]
		fields[fieldName] = fieldValue
	}

	// Create the entry
	entry := &Entry{
		ID:     entryID,
		Fields: fields,
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Create stream if it doesn't exist
	if _, exists := s.storage[key]; !exists {
		s.storage[key] = &Stream{
			entries: make([]*Entry, 0),
		}
	}

	// Append entry to the stream
	s.storage[key].entries = append(s.storage[key].entries, entry)

	// Return the entry ID as a bulk string
	return resp.MakeBulkString(entryID)
}
