package string_commands

import (
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

// Set sets the string value of a key.
// Example: SET mykey "Hello"
func (s *Store) Set(args []string) string {
	// SET command requires at least a key and a value
	if len(args) < 3 {
		return resp.MakeError("ERR wrong number of arguments for 'set' command")
	}

	key := args[1]
	value := args[2]

	var expiryType string
	var expiryValue int64
	expiryValue = 0
	if len(args) >= 5 {
		expiryType = strings.ToUpper(args[3])
		parsedValue, err := strconv.ParseInt(args[4], 10, 64)
		if err != nil {
			return resp.MakeError("ERR value is not an integer or out of range")
		}
		expiryValue = parsedValue

		if expiryType != "EX" && expiryType != "PX" {
			return resp.MakeError("ERR syntax error")
		}

		if expiryType == "EX" {
			expiryValue *= 1000
		}
	}

	// Store the key-value pair
	var expiryMilliseconds int64
	expiryMilliseconds = 0
	if expiryValue != 0 {
		expiryMilliseconds = time.Now().UnixMilli() + expiryValue
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.storage[key] = &StorageItem{
		Value:  value,
		Expiry: expiryMilliseconds,
	}

	// Return OK as a RESP simple string
	return resp.MakeSimpleString("OK")
}
