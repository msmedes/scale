package scale

import (
	"errors"
	"log"
	"sync"
	"time"
)

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
	a datastore
	a list of connections
	*/
}

type RemoteNode struct {
}

func (node *Node) join(other *RemoteNode) error {
	return errors.New("not implemented")
}

func (node *Node) stabilize(ticker *time.Ticker) {
	log.Fatal("not implemented")
}

func (node *Node) notify(remoteNode *RemoteNode) {
	log.Fatal("not implemented")
}

func (node *Node) findSuccessor(id []byte) (*RemoteNode, error) {
	return nil, errors.New("not implemented")
}

func (node *Node) findPredecessor(id []byte) (*RemoteNode, error) {
	return nil, errors.New("not implemented")
}
