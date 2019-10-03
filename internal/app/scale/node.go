package scale

import (
	"errors"
	"log"
	"time"

	uuid "github.com/google/uuid"
)

type Node struct {
	ID          Key
	predecessor *Node
	successor   *Node
	fingerTable FingerTable
	store       *Store
}

type RemoteNode struct {
	ID Key
}

func NewNode() *Node {
	n := &Node{
		ID:    genID(),
		store: NewStore(),
	}
	n.fingerTable = NewFingerTable(M, n)
	return n
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

func (node *Node) findSuccessor(ID []byte) (*RemoteNode, error) {
	return nil, errors.New("not implemented")
}

func (node *Node) findPredecessor(ID []byte) (*RemoteNode, error) {
	return nil, errors.New("not implemented")
}

func genID() Key {
	ID, err := uuid.NewRandom()

	if err != nil {
		log.Fatal(err)
	}

	return GenerateKey(ID.String())
}
