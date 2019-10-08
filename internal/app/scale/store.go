package scale

type data = map[Key][]byte

// Store underlying node KV store
type Store struct {
	data data
}

// NewStore set up a new store
func NewStore() *Store {
	return &Store{
		data: make(data),
	}
}

// Get get the given key
func (s *Store) Get(key Key) []byte {
	return s.data[key]
}

// Set set the given key
func (s *Store) Set(key Key, value []byte) error {
	s.data[key] = value
	return nil
}

// Del delete the given key
func (s *Store) Del(key Key) error {
	delete(s.data, key)
	return nil
}
