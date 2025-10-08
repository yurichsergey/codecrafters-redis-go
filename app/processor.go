package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type StorageItem struct {
	value  string
	expiry int64
}

type Processor struct {
	storage map[string]*StorageItem
}

func NewProcessor() *Processor {
	return &Processor{

		storage: make(map[string]*StorageItem),
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
	if command == "PING" {
		response = "+PONG\r\n"
	} else if command == "ECHO" {
		response = p.commandEcho(row)
	} else if command == "SET" {
		response = p.commandSet(row)
	} else if command == "GET" {
		response = p.commandGet(row)
	} else {
		response = "+PONG\r\n"
	}
	return response
}

func (p *Processor) commandEcho(strings []string) string {
	var response string
	response = "+"
	if len(strings) > 1 {
		response += strings[1:][0]
		for _, s := range strings[2:] {
			response += " " + s
		}
	}
	response += "\r\n"
	return response
}

func (p *Processor) commandSet(row []string) string {
	// SET command requires at least a key and a value
	if len(row) < 3 {
		return "-ERR wrong number of arguments for 'set' command\r\n"
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
			return "-ERR value is not an integer or out of range\r\n"
		}
		expiryValue = value

		if expiryType != "EX" && expiryType != "PX" {
			return "-ERR syntax error\r\n"
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
	p.storage[key] = &StorageItem{
		value:  value,
		expiry: expiryMilliseconds,
	}

	// Return OK as a RESP simple string
	return "+OK\r\n"
}

func (p *Processor) commandGet(row []string) string {
	// GET command requires a key argument
	if len(row) < 2 {
		return "-ERR wrong number of arguments for 'get' command\r\n"
	}

	key := row[1]

	// Check if the key exists in storage
	item, exists := p.storage[key]
	if !exists {
		// Return null bulk string if key doesn't exist
		return "$-1\r\n"
	}

	if item.expiry != 0 && time.Now().UnixMilli() > item.expiry {
		//delete(p.storage, key)
		return "$-1\r\n"
	}

	// Return the value as a RESP bulk string
	// Format: $<length>\r\n<data>\r\n
	return fmt.Sprintf("$%d\r\n%s\r\n", len(item.value), item.value)
}
