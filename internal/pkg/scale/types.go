package scale

// Canonical types for the implementation

// M bit keyspace
const M = 32

// Key 20 byte key
type Key = [M / 8]byte

// RemoteNode contains metadata (ID and Address) about another node in the network
type RemoteNode interface {
	GetAddr() string
	GetID() Key
	FindPredecessor(Key) (RemoteNode, error)
	FindSuccessor(Key) (RemoteNode, error)
	GetSuccessor() (RemoteNode, error)
	GetPredecessor() (RemoteNode, error)
	Ping() error
	Notify(Node) error
}

// Node represents the current node and operations it is responsible for
type Node interface {
	ClosestPrecedingFinger(Key) (RemoteNode, error)
	FindPredecessor(Key) (RemoteNode, error)
	FindSuccessor(Key) (RemoteNode, error)
	Get(Key) ([]byte, error)
	GetAddr() string
	GetFingerTableIDs() []Key
	GetID() Key
	GetLocal(Key) ([]byte, error)
	GetPort() string
	GetPredecessor() (RemoteNode, error)
	GetSuccessor() (RemoteNode, error)
	Notify(Key, string) error
	Set(Key, []byte) error
	SetLocal(Key, []byte) error
	TransferKeys(Key, string)
}

// Store represents a Scale-compatible underlying data store
type Store interface {
	Del(Key) error
	Get(Key) []byte
	Keys() []Key
	Set(Key, []byte) error
}
