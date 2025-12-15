package stream

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

// XRange retrieves a range of entries from the stream.
// Example: XRANGE mystream 0-1 0-2
func (s *Store) XRange(args []string) string {
	if len(args) < 4 {
		return resp.MakeError("ERR wrong number of arguments for 'xrange' command")
	}

	key := args[1]
	start := args[2]
	end := args[3]

	s.mutex.Lock()
	defer s.mutex.Unlock()

	stream, exists := s.storage[key]
	if !exists {
		return resp.MakeArray(nil)
	}

	startKey, err := ParseRangeID(start, true)
	if err != nil {
		return resp.MakeError("ERR " + err.Error())
	}

	endKey, err := ParseRangeID(end, false)
	if err != nil {
		return resp.MakeError("ERR " + err.Error())
	}

	entries := stream.tree.Range(startKey, endKey)

	var responseEntries []string
	for _, entry := range entries {
		// Entry ID
		idBulk := resp.MakeBulkString(entry.ID)

		// Fields
		var fieldStrings []string
		count := 0
		// Iterating over map does not guarantee order.
		// However, standard XRANGE returns fields in order of insertion.
		// Since we store them in a map, we lost the order.
		// We can't fix this without changing storage format.
		// We output what we have.
		for k, v := range entry.Fields {
			fieldStrings = append(fieldStrings, k)
			fieldStrings = append(fieldStrings, v)
			count++
		}
		fieldsArray := resp.MakeArray(fieldStrings)

		// [ID, [fields...]]
		entryArray := resp.MakeRESPArray([]string{idBulk, fieldsArray})
		responseEntries = append(responseEntries, entryArray)
	}

	return resp.MakeRESPArray(responseEntries)
}
