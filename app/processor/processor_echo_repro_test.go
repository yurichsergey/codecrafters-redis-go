package processor

import (
	"testing"
)

func TestEchoRepro(t *testing.T) {
	processor := NewProcessor()
	input := []string{"ECHO", "mango"}
	expected := "$5\r\nmango\r\n"
	result := processor.ProcessCommand(input)
	if result != expected {
		t.Errorf("ProcessCommand(%v) = %q, want %q", input, result, expected)
	}
}
