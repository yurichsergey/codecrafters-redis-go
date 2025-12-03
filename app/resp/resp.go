package resp

import (
	"fmt"
	"strings"
)

// MakeError creates a RESP error message.
func MakeError(msg string) string {
	return fmt.Sprintf("-%s\r\n", msg)
}

// MakeSimpleString creates a RESP simple string.
func MakeSimpleString(s string) string {
	return fmt.Sprintf("+%s\r\n", s)
}

// MakeBulkString creates a RESP bulk string.
func MakeBulkString(s string) string {
	return fmt.Sprintf("$%d\r\n%s\r\n", len(s), s)
}

// MakeInteger creates a RESP integer.
func MakeInteger(n int) string {
	return fmt.Sprintf(":%d\r\n", n)
}

// MakeNullBulkString creates a RESP null bulk string.
func MakeNullBulkString() string {
	return "$-1\r\n"
}

// MakeNullArray creates a RESP null array.
func MakeNullArray() string {
	return "*-1\r\n"
}

// MakeEmptyArray creates a RESP empty array.
func MakeEmptyArray() string {
	return "*0\r\n"
}

// MakeArray creates a RESP array of bulk strings.
func MakeArray(items []string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("*%d\r\n", len(items)))
	for _, item := range items {
		sb.WriteString(MakeBulkString(item))
	}
	return sb.String()
}
