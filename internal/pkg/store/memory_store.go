package store

// M keyspace
const M = 32

// Key 20 byte key
type Key = [M / 8]byte

type data = map[Key][]byte

// MemoryStore underlying node KV store
type MemoryStore struct {
	data data
}

// NewMemoryStore set up a new store
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(data),
	}
}

// Get get the given key
func (s *MemoryStore) Get(key Key) []byte {
	return s.data[key]
}

// Set set the given key
func (s *MemoryStore) Set(key Key, value []byte) error {
	s.data[key] = value
	return nil
}

// Del delete the given key
func (s *MemoryStore) Del(key Key) error {
	delete(s.data, key)
	return nil
}
