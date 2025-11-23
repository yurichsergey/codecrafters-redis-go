package main

import (
	"strings"
	"sync"
)

type StorageItem struct {
	value  string
	expiry int64
}

type BlockingResult struct {
	key   string
	value string
}

type BlockingClient struct {
	waiting chan BlockingResult
}

type Processor struct {
	storage         map[string]*StorageItem
	storageList     map[string][]string
	blockingClients map[string][]*BlockingClient
	clientsMutex    sync.Mutex
}

func NewProcessor() *Processor {
	return &Processor{
		storage:         make(map[string]*StorageItem),
		storageList:     make(map[string][]string),
		blockingClients: make(map[string][]*BlockingClient),
	}
}

func (p *Processor) ProcessCommand(row []string) string {
	var response string
	response = ""
	if len(row) == 0 {
		response = "$-1\r\n"
		return response
	}

	command := strings.ToUpper(row[0])
	switch command {
	case "PING":
		response = "+PONG\r\n"
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
		response = "+PONG\r\n"
	}
	return response
}
