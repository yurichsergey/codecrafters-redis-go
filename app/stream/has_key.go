package stream

// HasKey checks if a stream exists at the given key.
// Example: HasKey("mystream") returns true if "mystream" exists
func (s *Store) HasKey(key string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	_, exists := s.storage[key]
	return exists
}
