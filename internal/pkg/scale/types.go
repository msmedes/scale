package scale

// Canonical types for the implementation

// M bit keyspace
const M = 32

// Key 20 byte key
type Key = [M / 8]byte

type baseNode interface {
	GetAddr() string
	GetID() Key
	FindPredecessor(Key) (RemoteNode, error)
	FindSuccessor(Key) (RemoteNode, error)
	GetSuccessor() (RemoteNode, error)
	GetPredecessor() (RemoteNode, error)
	GetLocal(Key) ([]byte, error)
	SetLocal(Key, []byte) error
	ClosestPrecedingFinger(Key) (RemoteNode, error)
}

// RemoteNode contains metadata (ID and Address) about another node in the network
type RemoteNode interface {
	baseNode

	Ping() error
	Notify(Node) error
}

// Node represents the current node and operations it is responsible for
type Node interface {
	baseNode

	Get(Key) ([]byte, error)
	GetFingerTableIDs() []Key
	GetPort() string
	Notify(Key, string) error
	Set(Key, []byte) error
	TransferKeys(Key, string)
}

// Store represents a Scale-compatible underlying data store
type Store interface {
	Del(Key) error
	Get(Key) []byte
	Keys() []Key
	Set(Key, []byte) error
}
