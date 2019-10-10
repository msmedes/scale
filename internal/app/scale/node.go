package scale

import (
	"bytes"
	"context"
	"errors"
	"time"

	pb "github.com/msmedes/scale/internal/app/scale/proto"
	"go.uber.org/zap"
)

// StabilizeInterval how often to execute stabilization
const StabilizeInterval = 10

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

// StabilizationStart run a process that periodically makes sure the finger table
// is up to date and accurate
func (node *Node) StabilizationStart() {
	next := 0
	ticker := time.NewTicker(StabilizeInterval * time.Second)

	for {
		select {
		case <-ticker.C:
			next = node.fixNextFinger(next)
			node.stabilize()

			if next >= M {
				next = 0
			}
		}
	}
}

// Join join an existing network via another node
func (node *Node) Join(addr string) {
	node.logger.Infof("joining network via node at %s", addr)

	// create a client for the node we are trying to join
	remoteNode := NewRemoteNode(addr, node)
	node.remoteConnections[remoteNode.ID] = remoteNode

	// search for the successor to this node
	successor, err := remoteNode.RPC.FindSuccessor(
		context.Background(),
		&pb.RemoteQuery{Id: node.ID[:]},
	)

	successorID := ByteArrayToKey(successor.Id)

	_, ok := node.remoteConnections[successorID]

	// if the successor is not the node we are joining add it to remoteConnections
	if !ok {
		remoteNode = &RemoteNode{ID: successorID, Addr: successor.Addr, RPC: GetScaleClient(successor.Addr, node)}
		node.remoteConnections[successorID] = remoteNode
	}

	if err != nil {
		node.logger.Fatal(err)
	}

	node.successor = remoteNode
	node.logger.Infof("found successor: %s", KeyToString(node.successor.ID))
	node.logger.Info("joined network")
}

// FindSuccessor returns the successor for this node
func (node *Node) FindSuccessor(id Key) *RemoteNode {
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

// GetPredecessor returns the node's predecessor
func (node *Node) GetPredecessor() (*RemoteNode, error) {
	if node.predecessor != nil {
		return node.predecessor, nil
	}

	return nil, nil
}

// GetSuccessor retunrs the node's successor
func (node *Node) GetSuccessor() (*RemoteNode, error) {
	if node.successor != nil {
		return node.successor, nil
	}

	return nil, errors.New("no successor found")
}

// Shutdown leave the network
func (node *Node) Shutdown() {
	node.logger.Info("exiting")
}

func (node *Node) stabilize() {
	if bytes.Equal(node.successor.ID[:], node.ID[:]) {
		return
	}

	succPredecessor, err := node.successor.RPC.GetPredecessor(context.Background(), &pb.Empty{})

	if err != nil {
		node.logger.Error(err)
		node.logger.Fatalf("error retrieving predecessor from successor %s", KeyToString(node.successor.ID))
	}

	// The successor may not yet have a predecessor, meaning it has not
	// yet had a chance to update it's predecessor.  In that case we
	// notify the successor that we believe we are its predecessor.
	if succPredecessor.Present && Between(ByteArrayToKey(succPredecessor.Id), node.ID, node.successor.ID) {
		node.successor = NewRemoteNode(succPredecessor.Addr, node)
		node.logger.Infof("successor set to %s", KeyToString(node.successor.ID))
	}

	// tell the successor that node is the predecessor now
	node.successor.RPC.Notify(context.Background(), &pb.RemoteNode{Id: node.ID[:], Addr: node.Addr})
}

// Notify is called when another node thinks it is our predecessor
func (node *Node) Notify(id Key, addr string) error {
	if node.predecessor == nil || Between(id, node.predecessor.ID, node.ID) {
		node.predecessor = NewRemoteNode(addr, node)
		// Then we transfer keys
		node.logger.Infof("predecessor switched to %s", KeyToString(id))

		return nil
	}

	node.logger.Info("predecessor is already set")

	return nil
}

func (node *Node) fixNextFinger(next int) int {
	nextHash := fingerMath(node.ID[:], next, M)
	successor := node.FindSuccessor(ByteArrayToKey(nextHash))
	finger := &finger{ID: successor.ID}
	node.fingerTable[next] = finger
	return next + 1
}
