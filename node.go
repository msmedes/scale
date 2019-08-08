package scale

import "sync"

type Config struct {
}

// These are a bunch of guesses for what fields we'll need in the node struct
// We will likely need a mutex for each field
type Node struct {
	config *Config

	ID      []byte
	Address string

	Predecessor      *Node
	PredecessorMutex sync.RWMutex

	FingerTable      FingerTable
	FingerTableMutex syn.RWMutex
}
