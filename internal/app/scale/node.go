package scale

import (
	"bytes"
	"context"
	"time"

	pb "github.com/msmedes/scale/internal/app/scale/proto"
	"go.uber.org/zap"
)

// Node main node class
type Node struct {
	ID                Key
	Addr              string
	predecessor       *RemoteNode
	successor         *RemoteNode
	fingerTable       FingerTable
	store             *Store
	logger            *zap.SugaredLogger
	remoteConnections map[Key]*RemoteNode
}

// // NewRemoteNode create a new remote node with an RPC client
// func (node *Node) NewRemoteNode(id Key, addr string) *RemoteNode {
// 	return &RemoteNode{ID: id, Addr: addr, RPC: node.GetScaleClient(addr)}
// }

// NewNode create a new node
func NewNode(addr string, logger *zap.SugaredLogger) *Node {
	node := &Node{
		ID:                GenerateKey(addr),
		Addr:              addr,
		store:             NewStore(),
		logger:            logger,
		remoteConnections: make(map[Key]*RemoteNode),
	}

	node.fingerTable = NewFingerTable(M, node.ID)
	node.successor = &RemoteNode{ID: node.ID, Addr: node.Addr}

	return node
}

// Join join an existing network via another node
func (node *Node) Join(addr string) {
	node.logger.Infof("joining network via node at %s", addr)

	remoteNode := NewRemoteNode(addr, node)
	node.remoteConnections[remoteNode.ID] = remoteNode

	successor, err := remoteNode.RPC.FindSuccessor(
		context.Background(),
		&pb.RemoteQuery{Id: node.ID[:]},
	)

	successorID := ByteArrayToKey(successor.Id)

	_, ok := node.remoteConnections[successorID]

	if !ok {
		node.remoteConnections[successorID] = &RemoteNode{ID: successorID, Addr: successor.Addr, RPC: GetScaleClient(successor.Addr, node)}
	}

	if err != nil {
		node.logger.Fatal(err)
	}

	node.logger.Infof("found successor: %s", KeyToString(successorID))
	node.logger.Info("joined network")
}

// Shutdown leave the network
func (node *Node) Shutdown() {
	node.logger.Info("exiting")
}

func (node *Node) stabilize(ticker *time.Ticker) {
	node.logger.Fatal("not implemented")
}

func (node *Node) notify(remoteNode *RemoteNode) {
	node.logger.Fatal("not implemented")
}

// FindSuccessor returns the successor for this node
func (node *Node) FindSuccessor(id Key) *RemoteNode {
	return node.findSuccessor(id)

}

func (node *Node) findSuccessor(id Key) *RemoteNode {
	if BetweenRightInclusive(id, node.ID, node.successor.ID) {
		return node.successor
	}

	closestPrecedingID := node.closestPrecedingNode(id)
	if bytes.Equal(closestPrecedingID[:], node.ID[:]) {
		return &RemoteNode{ID: node.ID, Addr: node.Addr}
	}

	remoteNode, ok := node.remoteConnections[closestPrecedingID]
	node.logger.Infof("%+v", node.remoteConnections)

	if !ok {
		node.logger.Fatalf("remoteNode with ID %s not found", KeyToString(closestPrecedingID))
	}

	successor, err := remoteNode.RPC.FindSuccessor(
		context.Background(),
		&pb.RemoteQuery{Id: id[:]},
	)
	if err != nil {
		node.logger.Fatal(err)
	}

	return &RemoteNode{ID: ByteArrayToKey(successor.Id), Addr: successor.Addr}
}

// closestPrecedingNode returns the node in the finger table
// that is...the closest preceding node in the circle
func (node *Node) closestPrecedingNode(id Key) Key {
	// I think this could be implemented as binary search?
	for i := M - 1; i >= 0; i-- {
		finger := node.fingerTable[i]

		if Between(finger.ID, node.ID, id) {
			return finger.ID
		}
	}
	return node.ID
}

// FindPredecessor find predecessor
func (node *Node) FindPredecessor(ID Key) *RemoteNode {
	return node.predecessor
}
