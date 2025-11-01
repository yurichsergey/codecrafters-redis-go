package main

import (
	"sync"
	"testing"
	"time"
)

func TestBLPOPBasicFunctionality(t *testing.T) {
	tests := []struct {
		name     string
		setup    [][]string
		input    []string
		expected string
	}{
		{
			name:     "BLPOP on empty list",
			input:    []string{"BLPOP", "list_key", "0"},
			expected: "", // Will be blocked
		},
		{
			name:     "BLPOP with non-zero timeout (currently unsupported)",
			input:    []string{"BLPOP", "list_key", "10"},
			expected: "-ERR only timeout of 0 is currently supported\r\n",
		},
		{
			name:     "BLPOP without enough arguments",
			input:    []string{"BLPOP", "list_key"},
			expected: "-ERR wrong number of arguments for 'blpop' command\r\n",
		},
	}

	processor := NewProcessor()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processor.ProcessCommand(tt.input)
			if result != tt.expected {
				t.Errorf("ProcessCommand(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestBLPOPBlockingBehavior(t *testing.T) {
	t.Run("BLPOP blocks and receives element when pushed", func(t *testing.T) {
		processor := NewProcessor()

		// Use a WaitGroup to synchronize goroutines
		var wg sync.WaitGroup
		wg.Add(2)

		// Result channel to capture the BLPOP response
		resultChan := make(chan string, 1)

		// Goroutine for BLPOP (blocking call)
		go func() {
			defer wg.Done()
			blpopResult := processor.ProcessCommand([]string{"BLPOP", "block_list", "0"})
			resultChan <- blpopResult
		}()

		// Goroutine for RPUSH (to unblock)
		go func() {
			defer wg.Done()
			time.Sleep(50 * time.Millisecond) // Small delay to ensure BLPOP is blocking
			processor.ProcessCommand([]string{"RPUSH", "block_list", "element"})
		}()

		// Wait for both goroutines to complete
		wg.Wait()

		// Check the result
		select {
		case result := <-resultChan:
			expected := "*2\r\n$10\r\nblock_list\r\n$7\r\nelement\r\n"
			if result != expected {
				t.Errorf("BLPOP result = %q, want %q", result, expected)
			}
		case <-time.After(1 * time.Second):
			t.Fatal("BLPOP did not unblock within expected time")
		}
	})

	t.Run("BLPOP with multiple blocking clients", func(t *testing.T) {
		processor := NewProcessor()

		var wg sync.WaitGroup
		wg.Add(3)

		// Result channels for multiple blocking clients
		results := make(chan string, 2)

		// First blocking client
		go func() {
			defer wg.Done()
			result := processor.ProcessCommand([]string{"BLPOP", "multi_list", "0"})
			results <- result
		}()

		// Second blocking client
		go func() {
			defer wg.Done()
			result := processor.ProcessCommand([]string{"BLPOP", "multi_list", "0"})
			results <- result
		}()

		// RPUSH to unblock clients
		go func() {
			defer wg.Done()
			time.Sleep(50 * time.Millisecond)
			processor.ProcessCommand([]string{"RPUSH", "multi_list", "element"})
		}()

		// Wait for all goroutines to complete
		wg.Wait()

		// Collect results
		firstResult := <-results
		expectedFirst := "*2\r\n$10\r\nmulti_list\r\n$7\r\nelement\r\n"
		if firstResult != expectedFirst {
			t.Errorf("First BLPOP result = %q, want %q", firstResult, expectedFirst)
		}

		// Ensure the second client gets a different result or no result
		select {
		case secondResult := <-results:
			t.Errorf("Unexpected second BLPOP result: %q", secondResult)
		case <-time.After(100 * time.Millisecond):
			// Expected behavior - second client remains blocked
		}
	})
}

func TestBLPOPMultipleLists(t *testing.T) {
	t.Run("BLPOP with multiple lists", func(t *testing.T) {
		processor := NewProcessor()

		var wg sync.WaitGroup
		wg.Add(1)

		// Result channel to capture the BLPOP response
		resultChan := make(chan string, 1)

		// Goroutine for BLPOP on multiple lists
		go func() {
			defer wg.Done()
			blpopResult := processor.ProcessCommand([]string{"BLPOP", "list1", "list2", "list3", "0"})
			resultChan <- blpopResult
		}()

		// Goroutine for RPUSH to unblock
		go func() {
			time.Sleep(50 * time.Millisecond)
			processor.ProcessCommand([]string{"RPUSH", "list2", "element"})
		}()

		// Wait for goroutines to complete
		wg.Wait()

		// Check the result
		select {
		case result := <-resultChan:
			expected := "*2\r\n$5\r\nlist2\r\n$7\r\nelement\r\n"
			if result != expected {
				t.Errorf("BLPOP result = %q, want %q", result, expected)
			}
		case <-time.After(1 * time.Second):
			t.Fatal("BLPOP did not unblock within expected time")
		}
	})
}

func TestBLPOPCaseInsensitivity(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "BLPOP lowercase",
			input:    []string{"blpop", "list_key", "0"},
			expected: "", // Will be blocked
		},
		{
			name:     "BLPOP mixed case",
			input:    []string{"BlPop", "list_key", "0"},
			expected: "", // Will be blocked
		},
	}

	processor := NewProcessor()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processor.ProcessCommand(tt.input)
			if result != tt.expected {
				t.Errorf("ProcessCommand(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func BenchmarkBLPOPCommand(b *testing.B) {
	testCases := []struct {
		name  string
		setup []string
		input []string
	}{
		{
			name:  "BLPOP blocking scenario",
			setup: []string{"RPUSH", "bench_list", "element"},
			input: []string{"BLPOP", "bench_list", "0"},
		},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			processor := NewProcessor()

			// Setup the list for benchmarking
			processor.ProcessCommand(tc.setup)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				// In a real benchmark, you'd need to simulate the blocking/unblocking
				processor.ProcessCommand(tc.input)
			}
		})
	}
}
