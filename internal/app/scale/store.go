package scale

type data = map[Key][]byte

type Store struct {
	data data
}

func NewStore() *Store {
	return &Store{
		data: make(data),
	}
}

func (s *Store) Get(key Key) []byte {
	return s.data[key]
}

func (s *Store) Set(key Key, value []byte) error {
	s.data[key] = value
	return nil
}

func (s *Store) Del(key Key) error {
	delete(s.data, key)
	return nil
}
