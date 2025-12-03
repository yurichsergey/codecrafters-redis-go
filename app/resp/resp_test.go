package resp

import (
	"testing"
)

func TestMakeError(t *testing.T) {
	msg := "Error message"
	expected := "-Error message\r\n"
	result := MakeError(msg)
	if result != expected {
		t.Errorf("MakeError(%q) = %q; want %q", msg, result, expected)
	}
}

func TestMakeSimpleString(t *testing.T) {
	s := "OK"
	expected := "+OK\r\n"
	result := MakeSimpleString(s)
	if result != expected {
		t.Errorf("MakeSimpleString(%q) = %q; want %q", s, result, expected)
	}
}

func TestMakeBulkString(t *testing.T) {
	s := "Hello"
	expected := "$5\r\nHello\r\n"
	result := MakeBulkString(s)
	if result != expected {
		t.Errorf("MakeBulkString(%q) = %q; want %q", s, result, expected)
	}
}

func TestMakeInteger(t *testing.T) {
	n := 42
	expected := ":42\r\n"
	result := MakeInteger(n)
	if result != expected {
		t.Errorf("MakeInteger(%d) = %q; want %q", n, result, expected)
	}
}

func TestMakeNullBulkString(t *testing.T) {
	expected := "$-1\r\n"
	result := MakeNullBulkString()
	if result != expected {
		t.Errorf("MakeNullBulkString() = %q; want %q", result, expected)
	}
}

func TestMakeNullArray(t *testing.T) {
	expected := "*-1\r\n"
	result := MakeNullArray()
	if result != expected {
		t.Errorf("MakeNullArray() = %q; want %q", result, expected)
	}
}

func TestMakeEmptyArray(t *testing.T) {
	expected := "*0\r\n"
	result := MakeEmptyArray()
	if result != expected {
		t.Errorf("MakeEmptyArray() = %q; want %q", result, expected)
	}
}

func TestMakeArray(t *testing.T) {
	items := []string{"Hello", "World"}
	expected := "*2\r\n$5\r\nHello\r\n$5\r\nWorld\r\n"
	result := MakeArray(items)
	if result != expected {
		t.Errorf("MakeArray(%v) = %q; want %q", items, result, expected)
	}
}
