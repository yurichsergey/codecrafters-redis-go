package list_test

import "fmt"

// constructListResponse constructs a RESP array response from a slice of values.
// This helper is used by multiple test files to verify list contents.
func constructListResponse(values []string) string {
	// Construct the RESP array response
	var response string
	response = fmt.Sprintf("*%d\r\n", len(values))
	for _, item := range values {
		response += fmt.Sprintf("$%d\r\n%s\r\n", len(item), item)
	}

	return response
}
