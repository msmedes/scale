package node

import (
	"bytes"
	"context"
	"errors"
	"time"

	"github.com/msmedes/scale/internal/pkg/finger"
	"github.com/msmedes/scale/internal/pkg/keyspace"
	"github.com/msmedes/scale/internal/pkg/rpc"
	pb "github.com/msmedes/scale/internal/pkg/rpc/proto"
	"github.com/msmedes/scale/internal/pkg/scale"
	"github.com/msmedes/scale/internal/pkg/store"
	"go.uber.org/zap"
)

// StabilizeInterval how often to execute stabilization in seconds
const StabilizeInterval = 10

// Node main node class
type Node struct {
	scale.Node

	ID                keyspace.Key
	Addr              string
	predecessor       *RemoteNode
	successor         *RemoteNode
	fingerTable       finger.Table
	store             *store.MemoryStore
	Logger            *zap.SugaredLogger
	remoteConnections map[keyspace.Key]*RemoteNode
}

// NewNode create a new node
func NewNode(addr string) *Node {
	node := &Node{
		ID:                keyspace.GenerateKey(addr),
		Addr:              addr,
		store:             store.NewMemoryStore(),
		remoteConnections: make(map[keyspace.Key]*RemoteNode),
	}

	node.fingerTable = finger.NewFingerTable(keyspace.M, node.ID)
	node.successor = &RemoteNode{ID: node.ID, Addr: node.Addr}

	return node
}

// GetID getter for ID
func (node *Node) GetID() keyspace.Key {
	return node.ID
}

// GetAddr getter for address
func (node *Node) GetAddr() string {
	return node.Addr
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

			if next >= keyspace.M {
				next = 0
			}
		}
	}
}

// Join join an existing network via another node
func (node *Node) Join(addr string) {
	node.Logger.Infof("joining network via node at %s", addr)

	// create a client for the node we are trying to join
	remoteNode := NewRemoteNode(addr, node)
	node.remoteConnections[remoteNode.ID] = remoteNode

	// search for the successor to this node
	successor, err := remoteNode.RPC.FindSuccessor(
		context.Background(),
		&pb.RemoteQuery{Id: node.ID[:]},
	)

	successorID := keyspace.ByteArrayToKey(successor.Id)

	_, ok := node.remoteConnections[successorID]

	// if the successor is not the node we are joining add it to remoteConnections
	if !ok {
		remoteNode = &RemoteNode{
			ID:   successorID,
			Addr: successor.Addr,
			RPC:  rpc.NewClient(successor.Addr),
		}

		node.remoteConnections[successorID] = remoteNode
	}

	if err != nil {
		node.Logger.Fatal(err)
	}

	node.successor = remoteNode
	node.Logger.Infof("found successor: %s", keyspace.KeyToString(node.successor.ID))
	node.Logger.Info("joined network")
	node.stabilize()
	node.fixNextFinger(0)
}

// GetLocal return a value stored on this node
func (node *Node) GetLocal(key keyspace.Key) ([]byte, error) {
	return node.store.Get(key), nil
}

// SetLocal set a value in the local store
func (node *Node) SetLocal(key keyspace.Key, value []byte) error {
	return node.store.Set(key, value)
}

// Get return a value stored on this node
func (node *Node) Get(key keyspace.Key) ([]byte, error) {
	node.Logger.Infof("get: %x", key)
	if keyspace.BetweenRightInclusive(key, node.ID, node.successor.GetID()) {
		node.Logger.Infof("in between, getting on successor: %x", key)
		val, err := node.successor.RPC.GetLocal(
			context.Background(),
			&pb.GetRequest{Key: key[:]},
		)

		if err != nil {
			return nil, err
		}

		return val.GetValue(), nil
	}

	closestPrecedingID := node.closestPrecedingNode(key)

	if bytes.Equal(closestPrecedingID[:], node.ID[:]) {
		node.Logger.Infof("getting local: %x", key)
		return node.GetLocal(key)
	}

	remoteNode, ok := node.remoteConnections[closestPrecedingID]

	if !ok {
		node.Logger.Fatalf("remoteNode with ID %s not found", keyspace.KeyToString(closestPrecedingID))
	}

	node.Logger.Infof("calling get: %x on node %x", key, remoteNode.GetID())
	val, err := remoteNode.RPC.Get(
		context.Background(),
		&pb.GetRequest{Key: key[:]},
	)

	if err != nil {
		return nil, err
	}

	return val.GetValue(), nil
}

// Set set a value in the local store
func (node *Node) Set(key keyspace.Key, value []byte) error {
	node.Logger.Infof("set: %x", key)
	if keyspace.BetweenRightInclusive(key, node.ID, node.successor.GetID()) {
		node.Logger.Infof("between, setting on successor: %x", key)
		_, err := node.successor.RPC.SetLocal(context.Background(), &pb.SetRequest{
			Key:   key[:],
			Value: value,
		})

		if err != nil {
			return err
		}

		return nil
	}

	closestPrecedingID := node.closestPrecedingNode(key)

	if bytes.Equal(closestPrecedingID[:], node.ID[:]) {
		node.Logger.Infof("setting local: %x", key)
		return node.SetLocal(key, value)
	}

	remoteNode, ok := node.remoteConnections[closestPrecedingID]

	if !ok {
		node.Logger.Fatalf("remoteNode with ID %s not found", keyspace.KeyToString(closestPrecedingID))
	}

	node.Logger.Infof("calling set: %x on node %x", key, remoteNode.GetID())
	_, err := remoteNode.RPC.Set(
		context.Background(),
		&pb.SetRequest{Key: key[:], Value: value},
	)

	if err != nil {
		return err
	}

	return nil
}

// FindSuccessor returns the successor for this node
func (node *Node) FindSuccessor(id keyspace.Key) (scale.RemoteNode, error) {
	if keyspace.BetweenRightInclusive(id, node.ID, node.successor.GetID()) {
		return node.successor, nil
	}

	closestPrecedingID := node.closestPrecedingNode(id)

	if bytes.Equal(closestPrecedingID[:], node.ID[:]) {
		return &RemoteNode{ID: node.ID, Addr: node.Addr}, nil
	}

	remoteNode, ok := node.remoteConnections[closestPrecedingID]

	if !ok {
		node.Logger.Fatalf("remoteNode with ID %s not found", keyspace.KeyToString(closestPrecedingID))
	}

	successor, err := remoteNode.RPC.FindSuccessor(
		context.Background(),
		&pb.RemoteQuery{Id: id[:]},
	)

	if err != nil {
		node.Logger.Fatal(err)
	}

	return &RemoteNode{
		ID:   keyspace.ByteArrayToKey(successor.Id),
		Addr: successor.Addr,
	}, nil
}

// closestPrecedingNode returns the node in the finger table
// that is...the closest preceding node in the circle
func (node *Node) closestPrecedingNode(id keyspace.Key) keyspace.Key {
	// I think this could be implemented as binary search?
	for i := keyspace.M - 1; i >= 0; i-- {
		finger := node.fingerTable[i]

		if keyspace.Between(finger.ID, node.ID, id) {
			return finger.ID
		}
	}

	return node.ID
}

// GetPredecessor returns the node's predecessor
func (node *Node) GetPredecessor() (scale.RemoteNode, error) {
	if node.predecessor != nil {
		return node.predecessor, nil
	}

	return nil, nil
}

// GetSuccessor retunrs the node's successor
func (node *Node) GetSuccessor() (scale.RemoteNode, error) {
	if node.successor != nil {
		return node.successor, nil
	}

	return nil, errors.New("no successor found")
}

// Shutdown leave the network
func (node *Node) Shutdown() {
	node.Logger.Info("exiting")
}

func (node *Node) stabilize() {
	if bytes.Equal(node.successor.ID[:], node.ID[:]) {
		return
	}

	succPredecessor, err := node.successor.RPC.GetPredecessor(context.Background(), &pb.Empty{})

	if err != nil {
		node.Logger.Error(err)
		node.Logger.Fatalf("error retrieving predecessor from successor %s", keyspace.KeyToString(node.successor.ID))
	}

	// The successor may not yet have a predecessor, meaning it has not
	// yet had a chance to update it's predecessor.  In that case we
	// notify the successor that we believe we are its predecessor.
	if succPredecessor.Present && keyspace.Between(keyspace.ByteArrayToKey(succPredecessor.Id), node.ID, node.successor.ID) {
		node.successor = NewRemoteNode(succPredecessor.Addr, node)
		node.Logger.Infof("successor set to %s", keyspace.KeyToString(node.successor.ID))
	}

	// tell the successor that node is the predecessor now
	node.successor.RPC.Notify(context.Background(), &pb.RemoteNode{Id: node.ID[:], Addr: node.Addr})
}

// Notify is called when another node thinks it is our predecessor
func (node *Node) Notify(id keyspace.Key, addr string) error {
	if node.predecessor == nil || keyspace.Between(id, node.predecessor.ID, node.ID) {
		node.predecessor = NewRemoteNode(addr, node)
		// Then we transfer keys
		node.Logger.Infof("predecessor switched to %s", keyspace.KeyToString(id))

		return nil
	}

	node.Logger.Info("predecessor is already set")

	return nil
}

func (node *Node) fixNextFinger(next int) int {
	nextHash := finger.Math(node.ID[:], next, keyspace.M)
	successor, _ := node.FindSuccessor(keyspace.ByteArrayToKey(nextHash))
	finger := &finger.Finger{ID: successor.GetID()}
	node.fingerTable[next] = finger
	return next + 1
}
