package list

import (
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

// TestProcessor wraps Store to provide a ProcessCommand interface for testing
type TestProcessor struct {
	store *Store
}

// NewTestProcessor creates a new TestProcessor for testing
func NewTestProcessor() *TestProcessor {
	return &TestProcessor{
		store: NewStore(),
	}
}

// ProcessCommand handles the incoming Redis command and returns the response.
// This is a test helper that mimics the main Processor interface.
func (tp *TestProcessor) ProcessCommand(row []string) string {
	if len(row) == 0 {
		return resp.MakeNullBulkString()
	}

	command := strings.ToUpper(row[0])
	switch command {
	case "RPUSH":
		return tp.store.RPush(row)
	case "LRANGE":
		return tp.store.LRange(row)
	case "LPUSH":
		return tp.store.LPush(row)
	case "LLEN":
		return tp.store.LLen(row)
	case "LPOP":
		return tp.store.LPop(row)
	case "BLPOP":
		return tp.store.BLPop(row)
	default:
		return resp.MakeSimpleString("PONG")
	}
}
