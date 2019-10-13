package scale

import "github.com/msmedes/scale/internal/pkg/keyspace"

// Canonical types for the implementation

// RemoteNode contains metadata (ID and Address) about another node in the network
type RemoteNode interface {
	GetID() keyspace.Key
	GetAddr() string
}

// Node represents the current node and operations it is responsible for
type Node interface {
	Get(keyspace.Key) ([]byte, error)
	Set(keyspace.Key, []byte) error
	GetLocal(keyspace.Key) ([]byte, error)
	SetLocal(keyspace.Key, []byte) error
	Notify(keyspace.Key, string) error
	FindSuccessor(keyspace.Key) (RemoteNode, error)
	GetSuccessor() (RemoteNode, error)
	GetPredecessor() (RemoteNode, error)
	GetID() keyspace.Key
	GetAddr() string
	TransferKeys(keyspace.Key, string)
	GetFingerTableIDs() []keyspace.Key
}
