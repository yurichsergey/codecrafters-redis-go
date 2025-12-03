package processor

import (
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/list"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type StorageItem struct {
	// value is the string value stored in the item
	Value string
	// expiry is the expiration time in milliseconds
	Expiry int64
}

type Processor struct {
	// storage holds the key-value pairs for string commands
	Storage map[string]*StorageItem
	// listStore handles list-related commands
	ListStore *list.Store
}

// NewProcessor creates a new Processor instance with initialized storage and blocking clients.
func NewProcessor() *Processor {
	return &Processor{
		Storage:   make(map[string]*StorageItem),
		ListStore: list.NewStore(),
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
		response = p.CommandEcho(row)
	case "SET":
		response = p.CommandSet(row)
	case "GET":
		response = p.CommandGet(row)
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
