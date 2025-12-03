package list

// HasKey checks if a key exists in the list store.
// Returns true if the key exists, false otherwise.
func (s *Store) HasKey(key string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, exists := s.storage[key]
	return exists
}
