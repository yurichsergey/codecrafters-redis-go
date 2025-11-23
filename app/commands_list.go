package main

import (
	"fmt"
	"strconv"
	"strings"
)

// handleRPush appends one or more elements to the end of a list.
// Example: RPUSH mylist "hello" "world"
func (p *Processor) handleRPush(row []string) string {
	// Check if there are enough arguments
	if len(row) < 3 {
		return "-ERR wrong number of arguments for 'rpush' command\r\n"
	}

	key := row[1]
	elements := row[2:]

	p.storageList[key] = append(p.storageList[key], elements...)

	// Calculate the new length of the list
	newLength := len(p.storageList[key])

	// Wake up blocking clients for this key
	p.clientsMutex.Lock()
	defer p.clientsMutex.Unlock()

	if clients, exists := p.blockingClients[key]; exists {
		// Loop while we have both waiting clients and elements in the list
		for len(clients) > 0 && len(p.storageList[key]) > 0 {
			// Wake up the first (longest waiting) blocking client
			client := clients[0]
			client.waiting <- BlockingResult{key: key, value: p.storageList[key][0]}

			// Remove the first element and the first client
			p.storageList[key] = p.storageList[key][1:]
			clients = clients[1:]
		}

		// Update the blocking clients list
		p.blockingClients[key] = clients

		// Clean up empty list of blocking clients
		if len(p.blockingClients[key]) == 0 {
			delete(p.blockingClients, key)
		}
	}

	return fmt.Sprintf(":%d\r\n", newLength)
}

// handleLRange returns the specified elements of the list stored at key.
// Example: LRANGE mylist 0 -1
func (p *Processor) handleLRange(row []string) string {
	// Check if there are enough arguments
	if len(row) != 4 {
		return "-ERR wrong number of arguments for 'lrange' command\r\n"
	}

	key := row[1]
	startStr := row[2]
	stopStr := row[3]

	// Parse start index with negative index support
	start, err := strconv.Atoi(startStr)
	if err != nil {
		return "-ERR invalid start index\r\n"
	}

	// Parse stop index with negative index support
	stop, err := strconv.Atoi(stopStr)
	if err != nil {
		return "-ERR invalid stop index\r\n"
	}

	// Retrieve the list
	list, exists := p.storageList[key]
	if !exists {
		// If list doesn't exist, return an empty array
		return "*0\r\n"
	}

	// Calculate the length of the list
	listLength := len(list)

	// Handle negative indexes for start
	if start < 0 {
		start = listLength + start
		// If negative index is out of range, treat as 0
		if start < 0 {
			start = 0
		}
	}

	// Handle negative indexes for stop
	if stop < 0 {
		stop = listLength + stop
		// If negative index is out of range, treat as 0
		if stop < 0 {
			stop = 0
		}
	}

	// Adjust stop index if it's greater than or equal to list length
	if stop >= listLength {
		stop = listLength - 1
	}

	// Check if start index is out of bounds
	if start >= listLength {
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

// handleLPush inserts one or more elements at the head of a list.
// Example: LPUSH mylist "world"
func (p *Processor) handleLPush(row []string) string {
	// Check if there are enough arguments
	if len(row) < 3 {
		return "-ERR wrong number of arguments for 'lpush' command\r\n"
	}

	key := row[1]
	elements := row[2:]

	// Prepend elements in reverse order to maintain the correct list order
	for i, j := 0, len(elements)-1; i < j; i, j = i+1, j-1 {
		elements[i], elements[j] = elements[j], elements[i]
	}
	p.storageList[key] = append(elements, p.storageList[key]...)

	return fmt.Sprintf(":%d\r\n", len(p.storageList[key]))
}

// handleLLen returns the length of the list stored at key.
// Example: LLEN mylist
func (p *Processor) handleLLen(row []string) string {
	// Check if there are enough arguments
	if len(row) != 2 {
		return "-ERR wrong number of arguments for 'llen' command\r\n"
	}

	key := row[1]

	// Get the list length
	list, exists := p.storageList[key]
	if !exists {
		// If list doesn't exist, return 0
		return ":0\r\n"
	}

	// Return the length of the list as a RESP integer
	return fmt.Sprintf(":%d\r\n", len(list))
}

// handleLPop removes and returns the first elements of the list stored at key.
// Example: LPOP mylist
func (p *Processor) handleLPop(row []string) string {
	if len(row) < 2 {
		return "-ERR wrong number of arguments for 'lpop' command\r\n"
	}

	key := row[1]

	// Default count is 1
	count := 1

	// If count argument is provided
	if len(row) >= 3 {
		var err error
		count, err = strconv.Atoi(row[2])
		if err != nil || count < 0 {
			return "-ERR value is not an integer or out of range\r\n"
		}
	}

	// Check if list exists
	list, exists := p.storageList[key]
	if !exists || len(list) == 0 {
		return "$-1\r\n"
	}

	// Determine how many elements to actually remove
	numToRemove := count
	if numToRemove > len(list) {
		numToRemove = len(list)
	}

	// If count is 1 (no count argument provided), return single bulk string
	if len(row) == 2 {
		removed := list[0]
		p.storageList[key] = list[1:]

		// Clean up empty list
		if len(p.storageList[key]) == 0 {
			delete(p.storageList, key)
		}

		return fmt.Sprintf("$%d\r\n%s\r\n", len(removed), removed)
	}

	// Remove elements from the front
	removed := list[:numToRemove]
	p.storageList[key] = list[numToRemove:]

	// Clean up empty list
	if len(p.storageList[key]) == 0 {
		delete(p.storageList, key)
	}

	// Return as RESP array
	var response strings.Builder
	response.WriteString(fmt.Sprintf("*%d\r\n", len(removed)))
	for _, elem := range removed {
		response.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(elem), elem))
	}

	return response.String()
}

// handleBLPop removes and returns the first element of the list stored at key,
// blocking if the list is empty.
// Example: BLPOP mylist 0
func (p *Processor) handleBLPop(row []string) string {
	// Check if there are enough arguments
	if len(row) < 3 {
		return "-ERR wrong number of arguments for 'blpop' command\r\n"
	}

	// Currently only support timeout of 0
	timeout, err := strconv.Atoi(row[len(row)-1])
	if err != nil || timeout != 0 {
		return "-ERR only timeout of 0 is currently supported\r\n"
	}

	// Get the keys
	keys := row[1 : len(row)-1]

	// Check each list for an element
	for _, key := range keys {
		list, exists := p.storageList[key]
		if exists && len(list) > 0 {
			// Pop the first element
			element := list[0]
			p.storageList[key] = list[1:]

			// Clean up empty list
			if len(p.storageList[key]) == 0 {
				delete(p.storageList, key)
			}

			// Return the key and element as a RESP array
			return fmt.Sprintf("*2\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n",
				len(key), key, len(element), element)
		}
	}

	// If no elements are available, create a blocking client
	blockingClient := &BlockingClient{
		waiting: make(chan BlockingResult, 1),
	}

	p.clientsMutex.Lock()
	for _, key := range keys {
		p.blockingClients[key] = append(p.blockingClients[key], blockingClient)
	}
	p.clientsMutex.Unlock()

	// Block until an element is available
	result := <-blockingClient.waiting

	// Return the key and element
	return fmt.Sprintf("*2\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n",
		len(result.key), result.key, len(result.value), result.value)
}
