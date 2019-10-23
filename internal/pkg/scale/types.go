package scale

// Canonical types for the implementation

// M bit keyspace
const M = 32

// Key 20 byte key
type Key = [M / 8]byte

type baseNode interface {
	ClosestPrecedingFinger(Key) (RemoteNode, error)
	FindPredecessor(Key) (RemoteNode, error)
	FindSuccessor(Key) (RemoteNode, error)
	GetAddr() string
	GetID() Key
	GetLocal(Key) ([]byte, error)
	GetPredecessor() (RemoteNode, error)
	GetSuccessor() (RemoteNode, error)
	SetLocal(Key, []byte) error
	SetPredecessor(Key, string) error
	SetSuccessor(Key, string) error
}

// RemoteNode contains metadata (ID and Address) about another node in the network
type RemoteNode interface {
	baseNode

	CloseConnection() error
	Notify(Node) error
	Ping() error
}

// Node represents the current node and operations it is responsible for
type Node interface {
	baseNode

	Get(Key) ([]byte, error)
	GetFingerTableIDs() []Key
	GetPort() string
	GetKeys() []string
	Notify(Key, string) error
	Set(Key, []byte) error
	TransferKeys(Key, string) int
}

// Store represents a Scale-compatible underlying data store
type Store interface {
	Del(Key) error
	Get(Key) []byte
	Keys() []Key
	Set(Key, []byte) error
}
