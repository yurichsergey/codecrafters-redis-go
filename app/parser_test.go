package main

import (
	"reflect"
	"testing"
)

func TestParseString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
		hasError bool
	}{
		{
			name:     "Valid example from requirement",
			input:    "*4\r\n$2\r\nee\r\n$6\r\njhgjhg\r\n$3\r\nkjg\r\n$2\r\nkl\r\n",
			expected: []string{"ee", "jhgjhg", "kjg", "kl"},
			hasError: false,
		},
		{
			name:     "Single element",
			input:    "*1\r\n$5\r\nhello\r\n",
			expected: []string{"hello"},
			hasError: false,
		},
		{
			name:     "Two elements",
			input:    "*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n",
			expected: []string{"foo", "bar"},
			hasError: false,
		},
		{
			name:     "Empty strings",
			input:    "*2\r\n$0\r\n\r\n$3\r\nbar\r\n",
			expected: []string{"", "bar"},
			hasError: false,
		},
		{
			name:     "Zero elements",
			input:    "*0\r\n",
			expected: []string{},
			hasError: false,
		},
		{
			name:     "Single character strings",
			input:    "*3\r\n$1\r\na\r\n$1\r\nb\r\n$1\r\nc\r\n",
			expected: []string{"a", "b", "c"},
			hasError: false,
		},
		{
			name:     "Numbers as strings",
			input:    "*2\r\n$3\r\n123\r\n$1\r\n0\r\n",
			expected: []string{"123", "0"},
			hasError: false,
		},
		{
			name:     "Empty input",
			input:    "",
			expected: nil,
			hasError: true,
		},
		{
			name:     "Missing array indicator",
			input:    "4\r\n$2\r\nee\r\n",
			expected: nil,
			hasError: true,
		},
		{
			name:     "Invalid array count",
			input:    "*abc\r\n$2\r\nee\r\n",
			expected: nil,
			hasError: true,
		},
		{
			name:     "Missing bulk string indicator",
			input:    "*1\r\n2\r\nee\r\n",
			expected: nil,
			hasError: true,
		},
		{
			name:     "Invalid string length",
			input:    "*1\r\n$abc\r\nee\r\n",
			expected: nil,
			hasError: true,
		},
		{
			name:     "Length mismatch - too short",
			input:    "*1\r\n$5\r\nee\r\n",
			expected: nil,
			hasError: true,
		},
		{
			name:     "Length mismatch - too long",
			input:    "*1\r\n$1\r\nee\r\n",
			expected: nil,
			hasError: true,
		},
		{
			name:     "Missing string content",
			input:    "*1\r\n$2\r\n",
			expected: nil,
			hasError: true,
		},
		{
			name:     "Element count mismatch - too few",
			input:    "*3\r\n$2\r\nee\r\n$2\r\nkl\r\n",
			expected: nil,
			hasError: true,
		},
		{
			name:     "Special characters in strings",
			input:    "*2\r\n$3\r\n!@#\r\n$5\r\n$%^&*\r\n",
			expected: []string{"!@#", "$%^&*"},
			hasError: false,
		},
		{
			name:     "Spaces in strings",
			input:    "*2\r\n$11\r\nhello world\r\n$3\r\n   \r\n",
			expected: []string{"hello world", "   "},
			hasError: false,
		},
		{
			name:     "Large number of elements",
			input:    "*10\r\n$1\r\na\r\n$1\r\nb\r\n$1\r\nc\r\n$1\r\nd\r\n$1\r\ne\r\n$1\r\nf\r\n$1\r\ng\r\n$1\r\nh\r\n$1\r\ni\r\n$1\r\nj\r\n",
			expected: []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseString(tt.input)

			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestParseStringEdgeCases(t *testing.T) {
	t.Run("Input without final CRLF", func(t *testing.T) {
		input := "*1\r\n$5\r\nhello"
		result, err := parseString(input)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		expected := []string{"hello"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Large string", func(t *testing.T) {
		largeString := string(make([]byte, 1000))
		for i := range largeString {
			largeString = largeString[:i] + "a" + largeString[i+1:]
		}
		input := "*1\r\n$1000\r\n" + largeString + "\r\n"
		result, err := parseString(input)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		expected := []string{largeString}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected string of length %d, got length %d", len(expected[0]), len(result[0]))
		}
	})
}

func BenchmarkParseString(b *testing.B) {
	input := "*4\r\n$2\r\nee\r\n$6\r\njhgjhg\r\n$3\r\nkjg\r\n$2\r\nkl\r\n"

	for i := 0; i < b.N; i++ {
		_, err := parseString(input)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}
