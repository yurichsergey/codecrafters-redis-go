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
	storage     map[string]*StorageItem
	storageList map[string][]string
}

func NewProcessor() *Processor {
	return &Processor{
		storage:     make(map[string]*StorageItem),
		storageList: make(map[string][]string),
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
	default:
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
		// Return null bulk string if the key doesn't exist
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

// New function to handle RPUSH command
func (p *Processor) handleRPush(row []string) string {
	// Check if there are enough arguments
	if len(row) < 3 {
		return "-ERR wrong number of arguments for 'rpush' command\r\n"
	}

	key := row[1]
	elements := row[2:]

	p.storageList[key] = append(p.storageList[key], elements...)

	return fmt.Sprintf(":%d\r\n", len(p.storageList[key]))
}

// New function to handle LRANGE command
func (p *Processor) handleLRange(row []string) string {
	// Check if there are enough arguments
	if len(row) != 4 {
		return "-ERR wrong number of arguments for 'lrange' command\r\n"
	}

	key := row[1]
	startStr := row[2]
	stopStr := row[3]

	// Parse start and stop indexes
	start, err := strconv.Atoi(startStr)
	if err != nil || start < 0 {
		return "-ERR invalid start index\r\n"
	}

	stop, err := strconv.Atoi(stopStr)
	if err != nil || stop < 0 {
		return "-ERR invalid stop index\r\n"
	}

	// Retrieve the list
	list, exists := p.storageList[key]
	if !exists {
		// If list doesn't exist, return an empty array
		return "*0\r\n"
	}

	// Adjust stop index if it's greater than or equal to list length
	if stop >= len(list) {
		stop = len(list) - 1
	}

	// Check if start index is out of bounds
	if start >= len(list) {
		return "*0\r\n"
	}

	// Check if start index is greater than stop index
	if start > stop {
		return "*0\r\n"
	}

	// Extract the sublist
	subList := list[start : stop+1]

	// Construct the RESP array response
	var response string
	response = fmt.Sprintf("*%d\r\n", len(subList))
	for _, item := range subList {
		response += fmt.Sprintf("$%d\r\n%s\r\n", len(item), item)
	}

	return response
}
