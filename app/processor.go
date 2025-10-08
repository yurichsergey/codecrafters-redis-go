package main

import (
	"fmt"
	"strconv"
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

func (p *Processor) ProcessCommand(strings []string) string {
	var response string
	response = ""
	if len(strings) == 0 {
		response = ""
	} else if strings[0] == "PING" {
		response = "+PONG\r\n"
	} else if strings[0] == "ECHO" {
		response = p.commandEcho(strings)
	} else if strings[0] == "SET" {
		response = p.commandSet(strings)
	} else if strings[0] == "GET" {
		response = p.commandGet(strings)
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

func (p *Processor) commandSet(strings []string) string {
	// SET command requires at least a key and a value
	if len(strings) < 3 {
		return "-ERR wrong number of arguments for 'set' command\r\n"
	}

	key := strings[1]
	value := strings[2]

	var expiryType string
	var expiryValue int64
	expiryValue = 0
	if len(strings) >= 5 {
		expiryType = strings[3]
		value, err := strconv.ParseInt(strings[4], 10, 64)
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

func (p *Processor) commandGet(strings []string) string {
	// GET command requires a key argument
	if len(strings) < 2 {
		return "-ERR wrong number of arguments for 'get' command\r\n"
	}

	key := strings[1]

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
