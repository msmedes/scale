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

	FingerTable      FingerTable
	fingerTableMutex sync.RWMutex
}
