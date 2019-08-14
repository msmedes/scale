package scale

import "sync"

type Config struct {
}

type Node struct {
	config *Config

	ID      []byte
	Address string

	Predecessor      *Node
	predecessorMutex sync.RWMutex

	Successor      *Node
	successorMutex sync.RWMutex

	FingerTable      FingerTable
	fingerTableMutex sync.RWMutex

	/* eventually will probably need:
	a grpc server

	a datastore

	a list of connections

	*/
}
