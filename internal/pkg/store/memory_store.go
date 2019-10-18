package store

import (
	"sync"

	"github.com/msmedes/scale/internal/pkg/scale"
)

type data = map[scale.Key][]byte

// MemoryStore underlying node KV store
type MemoryStore struct {
	scale.Store

	mutex sync.RWMutex
	data  data
}

// NewMemoryStore set up a new store
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(data),
	}
}

// Get get the given key
func (s *MemoryStore) Get(key scale.Key) []byte {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.data[key]
}

// Set set the given key
func (s *MemoryStore) Set(key scale.Key, value []byte) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.data[key] = value
	return nil
}

// Del delete the given key
func (s *MemoryStore) Del(key scale.Key) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.data, key)
	return nil
}

// Keys list of keys
func (s *MemoryStore) Keys() []scale.Key {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	keys := make([]scale.Key, 0, len(s.data))

	for k := range s.data {
		keys = append(keys, k)
	}

	return keys
}
