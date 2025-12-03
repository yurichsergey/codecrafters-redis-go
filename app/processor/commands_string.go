package processor

import (
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

// CommandEcho returns the message passed to it.
// Example: ECHO "Hello World"
func (p *Processor) CommandEcho(strings []string) string {
	var content string
	if len(strings) > 1 {
		content += strings[1:][0]
		for _, s := range strings[2:] {
			content += " " + s
		}
	}
	return resp.MakeBulkString(content)
}

// CommandSet sets the string value of a key.
// Example: SET mykey "Hello"
func (p *Processor) CommandSet(row []string) string {
	// SET command requires at least a key and a value
	if len(row) < 3 {
		return resp.MakeError("ERR wrong number of arguments for 'set' command")
	}

	key := row[1]
	value := row[2]

	var expiryType string
	var expiryValue int64
	expiryValue = 0
	if len(row) >= 5 {
		expiryType = strings.ToUpper(row[3])
		value, err := strconv.ParseInt(row[4], 10, 64)
		if err != nil {
			return resp.MakeError("ERR value is not an integer or out of range")
		}
		expiryValue = value

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
	p.Storage[key] = &StorageItem{
		Value:  value,
		Expiry: expiryMilliseconds,
	}

	// Return OK as a RESP simple string
	return resp.MakeSimpleString("OK")
}

// CommandGet gets the value of a key.
// Example: GET mykey
func (p *Processor) CommandGet(row []string) string {
	// GET command requires a key argument
	if len(row) < 2 {
		return resp.MakeError("ERR wrong number of arguments for 'get' command")
	}

	key := row[1]

	// Check if the key exists in storage
	item, exists := p.Storage[key]
	if !exists {
		// Return null bulk string if the key doesn't exist
		return resp.MakeNullBulkString()
	}

	if item.Expiry != 0 && time.Now().UnixMilli() > item.Expiry {
		//delete(p.storage, key)
		return resp.MakeNullBulkString()
	}

	// Return the value as a RESP bulk string
	// Format: $<length>\r\n<data>\r\n
	return resp.MakeBulkString(item.Value)
}
