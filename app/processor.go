package main

import (
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/list"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type StorageItem struct {
	// value is the string value stored in the item
	value string
	// expiry is the expiration time in milliseconds
	expiry int64
}

type Processor struct {
	// storage holds the key-value pairs for string commands
	storage map[string]*StorageItem
	// listStore handles list-related commands
	listStore *list.Store
}

// NewProcessor creates a new Processor instance with initialized storage and blocking clients.
func NewProcessor() *Processor {
	return &Processor{
		storage:   make(map[string]*StorageItem),
		listStore: list.NewStore(),
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
		response = p.commandEcho(row)
	case "SET":
		response = p.commandSet(row)
	case "GET":
		response = p.commandGet(row)
	case "RPUSH":
		response = p.listStore.RPush(row)
	case "LRANGE":
		response = p.listStore.LRange(row)
	case "LPUSH":
		response = p.listStore.LPush(row)
	case "LLEN":
		response = p.listStore.LLen(row)
	case "LPOP":
		response = p.listStore.LPop(row)
	case "BLPOP":
		response = p.listStore.BLPop(row)
	default:
		response = resp.MakeSimpleString("PONG")
	}
	return response
}
