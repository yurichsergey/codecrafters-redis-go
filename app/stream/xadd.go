package stream

import (
	"fmt"
	"strconv"
	"strings"

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

	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Create stream if it doesn't exist
	stream, exists := s.storage[key]
	if !exists {
		stream = &Stream{
			tree: NewRadixTree(),
		}
		s.storage[key] = stream
	}

	// Get the last entry ID
	// Get the last entry ID
	lastID := ""
	if lastEntry := stream.tree.Last(); lastEntry != nil {
		lastID = lastEntry.ID
	}

	// Handle auto-generated sequence number logic: <time>-*
	if strings.HasSuffix(entryID, "-*") {
		timePart := strings.TrimSuffix(entryID, "-*")
		msTime, err := strconv.ParseInt(timePart, 10, 64)
		if err != nil {
			return resp.MakeError("ERR value is not an integer or out of range")
		}
		seqNum, err := GenerateSequence(msTime, lastID)
		if err != nil {
			return resp.MakeError(err.Error())
		}
		entryID = fmt.Sprintf("%d-%d", msTime, seqNum)
	}

	// Validate the new entry ID
	if err := ValidateID(entryID, lastID); err != nil {
		return resp.MakeError(err.Error())
	}

	// Create the entry
	entry := &Entry{
		ID:     entryID,
		Fields: fields,
	}

	// Append entry to the stream
	// Append entry to the stream
	keyStr, err := IDToKey(entryID)
	if err != nil {
		return resp.MakeError(err.Error())
	}
	stream.tree.Insert(keyStr, entry)

	// Return the entry ID as a bulk string
	return resp.MakeBulkString(entryID)
}
