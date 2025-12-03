package string_commands

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

// Echo returns the message passed to it.
// Example: ECHO "Hello World"
func (s *Store) Echo(args []string) string {
	var content string
	if len(args) > 1 {
		content += args[1:][0]
		for _, s := range args[2:] {
			content += " " + s
		}
	}
	return resp.MakeBulkString(content)
}
