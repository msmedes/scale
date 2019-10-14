package scale

// Canonical types for the implementation

// M bit keyspace
const M = 32

// Key 20 byte key
type Key = [M / 8]byte

// RemoteNode contains metadata (ID and Address) about another node in the network
type RemoteNode interface {
	GetID() Key
	GetAddr() string
}

// Node represents the current node and operations it is responsible for
type Node interface {
	Get(Key) ([]byte, error)
	Set(Key, []byte) error
	GetLocal(Key) ([]byte, error)
	SetLocal(Key, []byte) error
	Notify(Key, string) error
	FindSuccessor(Key) (RemoteNode, error)
	GetSuccessor() (RemoteNode, error)
	GetPredecessor() (RemoteNode, error)
	GetID() Key
	GetAddr() string
	GetPort() string
	TransferKeys(Key, string)
	GetFingerTableIDs() []Key
}

// Store represents a Scale-compatible underlying data store
type Store interface {
	Get(Key) []byte
	Set(Key, []byte) error
	Del(Key) error
	Keys() []Key
}
