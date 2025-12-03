package parser

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseString parses a Redis protocol string and returns a slice of parsed elements
func ParseString(input string) ([]string, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("empty input string")
	}

	// Split by \r\n to get individual lines
	lines := strings.Split(input, "\r\n")

	// Remove empty last element if string ends with \r\n
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	if len(lines) < 1 || !strings.HasPrefix(lines[0], "*") {
		return nil, fmt.Errorf("invalid format: expected array indicator '*'")
	}

	// Parse the number of elements from first line (e.g., "*4" -> 4)
	countStr := lines[0][1:] // Remove the '*' prefix
	expectedCount, err := strconv.Atoi(countStr)
	if err != nil {
		return nil, fmt.Errorf("invalid array count: %s", countStr)
	}

	result := make([]string, 0, expectedCount)
	i := 1 // Start from second line

	for len(result) < expectedCount && i < len(lines) {
		// Each element should start with $ followed by length
		if !strings.HasPrefix(lines[i], "$") {
			return nil, fmt.Errorf("expected bulk string indicator '$' at line %d", i)
		}

		// Parse the length (e.g., "$2" -> 2)
		lengthStr := lines[i][1:] // Remove the '$' prefix
		expectedLength, err := strconv.Atoi(lengthStr)
		if err != nil {
			return nil, fmt.Errorf("invalid string length: %s", lengthStr)
		}

		i++ // Move to the actual string content
		if i >= len(lines) {
			return nil, fmt.Errorf("missing string content for element %d", len(result))
		}

		content := lines[i]
		if len(content) != expectedLength {
			return nil, fmt.Errorf("string length mismatch: expected %d, got %d", expectedLength, len(content))
		}

		result = append(result, content)
		i++ // Move to next element
	}

	if len(result) != expectedCount {
		return nil, fmt.Errorf("element count mismatch: expected %d, got %d", expectedCount, len(result))
	}

	return result, nil
}
