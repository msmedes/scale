package store

import "github.com/msmedes/scale/internal/pkg/scale"

type data = map[scale.Key][]byte

// MemoryStore underlying node KV store
type MemoryStore struct {
	scale.Store

	data data
}

// NewMemoryStore set up a new store
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(data),
	}
}

// Get get the given key
func (s *MemoryStore) Get(key scale.Key) []byte {
	return s.data[key]
}

// Set set the given key
func (s *MemoryStore) Set(key scale.Key, value []byte) error {
	s.data[key] = value
	return nil
}

// Del delete the given key
func (s *MemoryStore) Del(key scale.Key) error {
	delete(s.data, key)
	return nil
}

// Keys list of keys
func (s *MemoryStore) Keys() []scale.Key {
	keys := make([]scale.Key, 0, len(s.data))

	for k := range s.data {
		keys = append(keys, k)
	}

	return keys
}
