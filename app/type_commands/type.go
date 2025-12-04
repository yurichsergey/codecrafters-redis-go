package type_commands

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

// Type returns the type of value stored at a key.
// Example: TYPE mykey
func (s *Store) Type(args []string) string {
	// TYPE command requires a key argument
	if len(args) < 2 {
		return resp.MakeError("ERR wrong number of arguments for 'type' command")
	}

	key := args[1]

	// Check if key exists in string storage
	if s.StringStore.HasKey(key) {
		return resp.MakeSimpleString("string")
	}

	// Check if key exists in list storage
	if s.ListStore.HasKey(key) {
		return resp.MakeSimpleString("list")
	}

	// Check if key exists in stream storage
	if s.StreamStore.HasKey(key) {
		return resp.MakeSimpleString("stream")
	}

	// Key doesn't exist in any storage
	return resp.MakeSimpleString("none")
}
