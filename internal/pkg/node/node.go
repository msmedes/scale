package node

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/msmedes/scale/internal/pkg/keyspace"
	pb "github.com/msmedes/scale/internal/pkg/rpc/proto"
	"github.com/msmedes/scale/internal/pkg/scale"
	"github.com/msmedes/scale/internal/pkg/store"
	"go.uber.org/zap"
)

// StabilizeInterval how often to execute stabilization in seconds
const StabilizeInterval = 1

// Node main node class
type Node struct {
	scale.Node
	sync.RWMutex

	ID                scale.Key
	Addr              string
	Port              string
	predecessor       *RemoteNode
	successor         *RemoteNode
	fingerTable       Table
	store             *store.MemoryStore
	Logger            *zap.SugaredLogger
	remoteConnections map[scale.Key]*RemoteNode
	shutdownChannel   chan struct{}
}

// NewNode create a new node
func NewNode(addr string) *Node {
	port := addr[strings.LastIndex(addr, ":")+1:]
	node := &Node{
		ID:                keyspace.GenerateKey(addr),
		Addr:              addr,
		Port:              port,
		store:             store.NewMemoryStore(),
		remoteConnections: make(map[scale.Key]*RemoteNode),
		shutdownChannel:   make(chan struct{}),
	}

	node.fingerTable = NewScaleFingerTable(node)
	node.successor = NewRemoteNode(node.Addr, node)
	node.predecessor = NewRemoteNode(node.Addr, node)

	return node
}

// GetID getter for ID
func (node *Node) GetID() scale.Key {
	return node.ID
}

// GetAddr getter for address
func (node *Node) GetAddr() string {
	return node.Addr
}

// GetPort getter for port
func (node *Node) GetPort() string {
	return node.Port
}

// GetFingerTableIDs return an array of IDs in the table
func (node *Node) GetFingerTableIDs() []scale.Key {
	node.RLock()
	defer node.RUnlock()
	var keys []scale.Key

	for _, k := range node.fingerTable {
		keys = append(keys, k.ID)
	}

	return keys
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
			node.checkPredecessor()

			if next == scale.M {
				next = 0
			}
		case <-node.shutdownChannel:
			ticker.Stop()
			return
		}
	}
}

// TransferKeys transfer keys to the given node
func (node *Node) TransferKeys(id scale.Key, addr string) {
	remote := NewRemoteNode(addr, node)

	if bytes.Compare(id[:], node.ID[:]) >= 0 {
		return
	}

	for _, k := range node.store.Keys() {
		if bytes.Compare(k[:], id[:]) >= 0 {
			break
		}

		node.transferKey(k, remote)
	}
}

func (node *Node) transferKey(key scale.Key, remote *RemoteNode) {
	val, err := node.GetLocal(key)

	if err != nil {
		node.Logger.Error(err)
	}

	req := &pb.SetRequest{Key: key[:], Value: val[:]}
	_, err = remote.RPC.SetLocal(context.Background(), req)

	if err != nil {
		node.Logger.Error(err)
	}

	node.store.Del(key)
}

// Join join an existing network via another node
// This implements the pseudocode for in fig. 5 of the new paper
// for concurrent joins
func (node *Node) Join(addr string) {
	var remoteNode *RemoteNode
	node.Logger.Infof("joining network via node at %s", addr)

	// create a client for the node we are trying to join
	remoteNode = NewRemoteNode(addr, node)

	// search for the successor to this node
	successor, err := remoteNode.RPC.FindSuccessor(
		context.Background(),
		&pb.RemoteQuery{Id: node.ID[:]},
	)

	successorID := keyspace.ByteArrayToKey(successor.Id)

	_, ok := node.remoteConnections[successorID]

	// if the successor is not the node we are joining add it to remoteConnections
	if !ok {
		remoteNode = NewRemoteNode(successor.Addr, node)
	}

	if err != nil {
		node.Logger.Fatal(err)
	}

	node.successor = remoteNode

	if !keyspace.Equal(node.ID, node.successor.ID) {
		node.successor.RPC.TransferKeys(context.Background(), &pb.KeyTransferRequest{
			Id:   node.ID[:],
			Addr: node.Addr,
		})
	}

	node.Logger.Infof("found successor: %s", keyspace.KeyToString(node.successor.ID))
	node.Logger.Info("joined network")
}

// GetLocal return a value stored on this node
func (node *Node) GetLocal(key scale.Key) ([]byte, error) {
	node.RLock()
	defer node.RUnlock()

	return node.store.Get(key), nil
}

// SetLocal set a value in the local store
func (node *Node) SetLocal(key scale.Key, value []byte) error {
	node.RLock()
	defer node.RUnlock()

	return node.store.Set(key, value)
}

// Get return a value stored on this node
func (node *Node) Get(key scale.Key) ([]byte, error) {
	node.RLock()
	defer node.RUnlock()

	succ, err := node.FindSuccessor(key)
	remoteNode := NewRemoteNode(succ.GetAddr(), node)

	val, err := remoteNode.RPC.GetLocal(
		context.Background(),
		&pb.GetRequest{Key: key[:]},
	)

	if err != nil {
		return nil, err
	}

	return val.GetValue(), nil
}

// Set set a value in the local store
func (node *Node) Set(key scale.Key, value []byte) error {
	node.RLock()
	defer node.RUnlock()

	succ, err := node.FindSuccessor(key)
	remoteNode := NewRemoteNode(succ.GetAddr(), node)

	_, err = remoteNode.RPC.SetLocal(
		context.Background(),
		&pb.SetRequest{Key: key[:], Value: value},
	)

	if err != nil {
		return err
	}

	return nil
}

//FindPredecessor finds the predecessor to the id
func (node *Node) FindPredecessor(id scale.Key) (scale.RemoteNode, error) {
	node.Lock()
	node.Unlock()

	var closestPrecedingRPC *pb.RemoteNode
	if keyspace.Equal(node.ID, node.successor.ID) {
		return NewRemoteNode(node.Addr, node), nil // this is the only node in the network
	}

	// I apologize for this Swift looking code but I don't see a better way at the moment
	closestPreceding := NewRemoteNode(node.Addr, node)
	predecessorSuccessor, _ := closestPreceding.RPC.GetSuccessor(context.Background(), &pb.Empty{})
	predecessorSuccessorRemote := NewRemoteNode(predecessorSuccessor.Addr, node)

	for !keyspace.BetweenRightInclusive(id, closestPreceding.ID, predecessorSuccessorRemote.ID) {
		closestPrecedingRPC, _ = closestPreceding.RPC.ClosestPrecedingFinger(context.Background(), &pb.RemoteQuery{Id: id[:]})
		closestPreceding = NewRemoteNode(closestPrecedingRPC.Addr, node)
		predecessorSuccessor, _ = closestPreceding.RPC.GetSuccessor(context.Background(), &pb.Empty{})
		predecessorSuccessorRemote = NewRemoteNode(predecessorSuccessor.Addr, node)

	}

	return closestPreceding, nil

}

// FindSuccessor returns the successor for this node
func (node *Node) FindSuccessor(id scale.Key) (scale.RemoteNode, error) {
	predecessor, _ := node.FindPredecessor(id)
	// idk this is the only way I know how to do this right now.  The compiler was mad that
	// that FindPredecessor technically returns an interface
	predecessorRemote := predecessor.(*RemoteNode)

	successorRPC, _ := predecessorRemote.RPC.GetSuccessor(context.Background(), &pb.Empty{})
	return NewRemoteNode(successorRPC.Addr, node), nil

}

// FindSuccessor returns the successor for this node
// func (node *Node) FindSuccessor(id scale.Key) (scale.RemoteNode, error) {
// 	if keyspace.BetweenRightInclusive(id, node.ID, node.successor.GetID()) {
// 		return node.successor, nil
// 	}

// 	closestPrecedingID := node.closestPrecedingNode(id)

// 	if bytes.Equal(closestPrecedingID[:], node.ID[:]) {
// 		return &RemoteNode{ID: node.ID, Addr: node.Addr}, nil
// 	}

// 	remoteNode, ok := node.remoteConnections[closestPrecedingID]

// 	if !ok {
// 		node.Logger.Fatalf("remoteNode with ID %s not found", keyspace.KeyToString(closestPrecedingID))
// 	}

// 	successor, err := remoteNode.RPC.FindSuccessor(
// 		context.Background(),
// 		&pb.RemoteQuery{Id: id[:]},
// 	)

// 	if err != nil {
// 		node.Logger.Fatal(err)
// 	}

// 	return NewRemoteNode(successor.Addr, node), nil
// }

// ClosestPrecedingFinger returns the closest preceding finger to the id
func (node *Node) ClosestPrecedingFinger(id scale.Key) (scale.RemoteNode, error) {
	node.RLock()
	defer node.RUnlock()

	// I think this could be implemented as binary search?
	for i := scale.M - 1; i >= 0; i-- {
		finger := node.fingerTable[i]

		if keyspace.Between(finger.ID, node.ID, id) {
			return finger, nil
		}
	}
	return NewRemoteNode(node.Addr, node), nil

}

// closestPrecedingNode returns the node in the finger table
// that is...the closest preceding node in the circle
// func (node *Node) closestPrecedingNode(id scale.Key) scale.Key {
// 	node.RLock()
// 	defer node.RUnlock()

// 	// I think this could be implemented as binary search?
// 	for i := scale.M - 1; i >= 0; i-- {
// 		finger := node.fingerTable[i]

// 		if keyspace.Between(finger.ID, node.ID, id) {
// 			return finger.ID
// 		}
// 	}

// 	return node.ID
// }

// GetPredecessor returns the node's predecessor
func (node *Node) GetPredecessor() (scale.RemoteNode, error) {
	node.RLock()
	defer node.RUnlock()

	if node.predecessor != nil {
		return node.predecessor, nil
	}

	return nil, nil
}

// GetSuccessor retunrs the node's successor
func (node *Node) GetSuccessor() (scale.RemoteNode, error) {
	node.RLock()
	defer node.RUnlock()

	if node.successor != nil {
		return node.successor, nil
	}

	return nil, errors.New("no successor found")
}

// Shutdown leave the network
func (node *Node) Shutdown() {
	node.RLock()
	defer node.RUnlock()

	close(node.shutdownChannel)

	if !keyspace.Equal(node.ID, node.successor.ID) {
	}
	for _, remoteConnection := range node.remoteConnections {
		remoteConnection.clientConnection.Close()
	}
}

func (node *Node) stabilize() {

	node.RLock()
	defer node.RUnlock()

	if !keyspace.Equal(node.successor.ID, node.ID) {
		succPredecessor, err := node.successor.RPC.GetPredecessor(context.Background(), &pb.Empty{})

		if err != nil {
			node.Logger.Errorf("error retrieving predecessor from successor %s", keyspace.KeyToString(node.successor.ID))
			node.Logger.Error(err)
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

		return
	}
	if node.predecessor != nil {
		node.successor = node.predecessor
		node.Logger.Infof("successor set to %s", keyspace.KeyToString(node.successor.ID))
	}
}

func (node *Node) checkPredecessor() {
	node.RLock()
	defer node.RUnlock()
	predecessor := node.predecessor

	if predecessor == nil {
		return
	}

	_, err := predecessor.RPC.Ping(context.Background(), &pb.Empty{})

	if err != nil {
		node.Logger.Infof("predecessor unresponsive. removing %s", keyspace.KeyToString(predecessor.GetID()))
		node.predecessor = nil
	}
}

// Notify is called when another node thinks it is our predecessor
func (node *Node) Notify(id scale.Key, addr string) error {
	node.RLock()
	defer node.RUnlock()
	if node.predecessor == nil || keyspace.Between(id, node.predecessor.ID, node.ID) {
		node.predecessor = NewRemoteNode(addr, node)
		node.predecessor.RPC.TransferKeys(context.Background(), &pb.KeyTransferRequest{
			Id:   node.ID[:],
			Addr: node.Addr,
		})
		node.Logger.Infof("predecessor set to %s", keyspace.KeyToString(id))

		return nil
	}

	return nil
}

func (node *Node) fixNextFinger(next int) int {
	node.RLock()
	defer node.RUnlock()

	nextHash := FingerMath(node.ID[:], next, scale.M)
	successor, _ := node.FindSuccessor(keyspace.ByteArrayToKey(nextHash))
	successorRemote := successor.(*RemoteNode)
	finger := NewRemoteNode(successorRemote.Addr, node)
	node.fingerTable[next] = finger
	return next + 1
}
