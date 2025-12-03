package main

import (
	"strings"
	"sync"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type StorageItem struct {
	// value is the string value stored in the item
	value string
	// expiry is the expiration time in milliseconds
	expiry int64
}

type BlockingResult struct {
	// key is the list key that the client was waiting for
	key string
	// value is the element popped from the list
	value string
}

type BlockingClient struct {
	// waiting is a channel that receives the result when an element is available
	waiting chan BlockingResult
}

type Processor struct {
	// storage holds the key-value pairs for string commands
	storage map[string]*StorageItem
	// storageList holds the key-value pairs for list commands
	storageList map[string][]string
	// blockingClients holds the list of clients waiting for elements on specific keys
	blockingClients map[string][]*BlockingClient
	// clientsMutex protects access to the blockingClients map
	clientsMutex sync.Mutex
}

// NewProcessor creates a new Processor instance with initialized storage and blocking clients.
func NewProcessor() *Processor {
	return &Processor{
		storage:         make(map[string]*StorageItem),
		storageList:     make(map[string][]string),
		blockingClients: make(map[string][]*BlockingClient),
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
		response = p.handleRPush(row)
	case "LRANGE":
		response = p.handleLRange(row)
	case "LPUSH":
		response = p.handleLPush(row)
	case "LLEN":
		response = p.handleLLen(row)
	case "LPOP":
		response = p.handleLPop(row)
	case "BLPOP":
		response = p.handleBLPop(row)
	default:
		response = resp.MakeSimpleString("PONG")
	}
	return response
}
