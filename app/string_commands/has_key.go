package string_commands

import (
	"time"
)

// HasKey checks if a key exists in the string store and is not expired.
// Returns true if the key exists and is valid, false otherwise.
func (s *Store) HasKey(key string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	item, exists := s.storage[key]
	if !exists {
		return false
	}

	// Check if the key has expired
	if item.Expiry != 0 && time.Now().UnixMilli() > item.Expiry {
		return false
	}

	return true
}
