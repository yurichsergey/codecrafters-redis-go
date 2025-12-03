package processor

import (
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/list"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/string_commands"
)

type Processor struct {
	// StringStore handles string-related commands
	StringStore *string_commands.Store
	// listStore handles list-related commands
	ListStore *list.Store
}

// NewProcessor creates a new Processor instance with initialized storage and blocking clients.
func NewProcessor() *Processor {
	return &Processor{
		StringStore: string_commands.NewStore(),
		ListStore:   list.NewStore(),
	}
}

// ProcessCommand handles the incoming Redis command and returns the response.
func (p *Processor) ProcessCommand(row []string) string {
	var response string
	response = ""
	if len(row) == 0 {
		response = resp.MakeNullBulkString()
		return response
	}

	command := strings.ToUpper(row[0])
	switch command {
	case "PING":
		response = resp.MakeSimpleString("PONG")
	case "ECHO":
		response = p.StringStore.Echo(row)
	case "SET":
		response = p.StringStore.Set(row)
	case "GET":
		response = p.StringStore.Get(row)
	case "RPUSH":
		response = p.ListStore.RPush(row)
	case "LRANGE":
		response = p.ListStore.LRange(row)
	case "LPUSH":
		response = p.ListStore.LPush(row)
	case "LLEN":
		response = p.ListStore.LLen(row)
	case "LPOP":
		response = p.ListStore.LPop(row)
	case "BLPOP":
		response = p.ListStore.BLPop(row)
	default:
		response = resp.MakeSimpleString("PONG")
	}
	return response
}
