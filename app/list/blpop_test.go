package list_test

import (
	"sync"
	"testing"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/processor"
)

func TestBLPOPBasicFunctionality(t *testing.T) {
	tests := []struct {
		name     string
		setup    [][]string
		input    []string
		expected string
	}{
		// Removed "BLPOP on empty list" case as it blocks indefinitely in this synchronous test runner

		{
			name:     "BLPOP without enough arguments",
			input:    []string{"BLPOP", "list_key"},
			expected: "-ERR wrong number of arguments for 'blpop' command\r\n",
		},
	}

	processor := processor.NewProcessor()
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
		processor := processor.NewProcessor()

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
		processor := processor.NewProcessor()

		var wg sync.WaitGroup
		wg.Add(2)

		// Result channels for multiple blocking clients
		results := make(chan string, 2)

		// Start two blocking clients
		for i := 0; i < 2; i++ {
			go func() {
				defer wg.Done()
				result := processor.ProcessCommand([]string{"BLPOP", "multi_list", "0"})
				results <- result
			}()
		}

		// Allow time for clients to block
		time.Sleep(50 * time.Millisecond)

		// RPUSH to unblock one client
		processor.ProcessCommand([]string{"RPUSH", "multi_list", "element"})

		// Collect first result
		select {
		case firstResult := <-results:
			expectedFirst := "*2\r\n$10\r\nmulti_list\r\n$7\r\nelement\r\n"
			if firstResult != expectedFirst {
				t.Errorf("First BLPOP result = %q, want %q", firstResult, expectedFirst)
			}
		case <-time.After(1 * time.Second):
			t.Fatal("Timeout waiting for first client to be unblocked")
		}

		// Ensure the second client gets a different result or no result
		select {
		case secondResult := <-results:
			t.Errorf("Unexpected second BLPOP result: %q", secondResult)
		case <-time.After(100 * time.Millisecond):
			// Expected behavior - second client remains blocked
		}

		// Unblock the second client so the test can finish
		processor.ProcessCommand([]string{"RPUSH", "multi_list", "cleanup"})

		// Wait for both clients to finish
		wg.Wait()
	})

	t.Run("BLPOP wakes up multiple clients with multiple elements", func(t *testing.T) {
		processor := processor.NewProcessor()

		var wg sync.WaitGroup
		wg.Add(2)

		results := make(chan string, 2)

		// Client 1
		go func() {
			defer wg.Done()
			results <- processor.ProcessCommand([]string{"BLPOP", "multi_wake", "0"})
		}()

		// Client 2
		go func() {
			defer wg.Done()
			results <- processor.ProcessCommand([]string{"BLPOP", "multi_wake", "0"})
		}()

		// RPUSH with 2 elements
		go func() {
			time.Sleep(50 * time.Millisecond)
			processor.ProcessCommand([]string{"RPUSH", "multi_wake", "val1", "val2"})
		}()

		wg.Wait()

		// Verify both clients got results
		for i := 0; i < 2; i++ {
			select {
			case res := <-results:
				if res == "" {
					t.Errorf("Received empty result")
				}
			case <-time.After(1 * time.Second):
				t.Fatal("Timeout waiting for clients to wake up")
			}
		}
	})
}

func TestBLPOPMultipleLists(t *testing.T) {
	t.Run("BLPOP with multiple lists", func(t *testing.T) {
		processor := processor.NewProcessor()

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
			input:    []string{"blpop", "case_list", "0"},
			expected: "*2\r\n$9\r\ncase_list\r\n$7\r\nelement\r\n",
		},
		{
			name:     "BLPOP mixed case",
			input:    []string{"BlPop", "case_list", "0"},
			expected: "*2\r\n$9\r\ncase_list\r\n$7\r\nelement\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor := processor.NewProcessor()

			// Result channel
			resultChan := make(chan string, 1)

			// Blocking call
			go func() {
				resultChan <- processor.ProcessCommand(tt.input)
			}()

			// Unblock
			go func() {
				time.Sleep(10 * time.Millisecond)
				processor.ProcessCommand([]string{"RPUSH", "case_list", "element"})
			}()

			// Check result
			select {
			case result := <-resultChan:
				if result != tt.expected {
					t.Errorf("ProcessCommand(%v) = %q, want %q", tt.input, result, tt.expected)
				}
			case <-time.After(1 * time.Second):
				t.Fatal("BLPOP did not unblock")
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
			processor := processor.NewProcessor()

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

func TestRPushReturnValueWithBlockedClients(t *testing.T) {
	t.Run("RPUSH returns correct length when client is blocked", func(t *testing.T) {
		processor := processor.NewProcessor()

		var wg sync.WaitGroup
		wg.Add(1)

		// Blocking client
		go func() {
			defer wg.Done()
			processor.ProcessCommand([]string{"BLPOP", "return_val_list", "0"})
		}()

		// Ensure client is blocked
		time.Sleep(50 * time.Millisecond)

		// RPUSH
		result := processor.ProcessCommand([]string{"RPUSH", "return_val_list", "element"})

		wg.Wait()

		// Expect 1 because we pushed 1 element, even if it was immediately popped
		expected := ":1\r\n"
		if result != expected {
			t.Errorf("RPUSH result = %q, want %q", result, expected)
		}
	})
}
