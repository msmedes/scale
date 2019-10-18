package node

import (
	"bytes"
	"context"
	"errors"
	"log"
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

	id                scale.Key
	addr              string
	port              string
	predecessor       *RemoteNode
	successor         *RemoteNode
	fingerTable       Table
	store             *store.MemoryStore
	logger            *zap.Logger
	sugar             *zap.SugaredLogger
	remoteConnections map[scale.Key]*RemoteNode
	shutdownChannel   chan struct{}
	mutex             sync.RWMutex
}

// NewNode create a new node
func NewNode(addr string) *Node {
	port := addr[strings.LastIndex(addr, ":")+1:]

	node := &Node{
		id:                keyspace.GenerateKey(addr),
		addr:              addr,
		port:              port,
		store:             store.NewMemoryStore(),
		remoteConnections: make(map[scale.Key]*RemoteNode),
		shutdownChannel:   make(chan struct{}),
	}

	node.fingerTable = NewScaleFingerTable(node)
	node.successor = NewRemoteNode(node.addr)
	node.predecessor = NewRemoteNode(node.addr)

	logger, err := zap.NewDevelopment(
		zap.Fields(
			zap.String("node", keyspace.KeyToString(node.id)),
		),
	)

	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}

	sugar := logger.Sugar()
	node.sugar = sugar
	node.logger = logger

	return node
}

// GetID getter for ID
func (node *Node) GetID() scale.Key {
	return node.id
}

// GetAddr getter for address
func (node *Node) GetAddr() string {
	return node.addr
}

// GetPort getter for port
func (node *Node) GetPort() string {
	return node.port
}

// GetFingerTableIDs return an array of IDs in the table
func (node *Node) GetFingerTableIDs() []scale.Key {
	node.mutex.RLock()
	defer node.mutex.RUnlock()
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
	remote := NewRemoteNode(addr)

	if bytes.Compare(id[:], node.id[:]) >= 0 {
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
		node.sugar.Error(err)
	}

	req := &pb.SetRequest{Key: key[:], Value: val[:]}
	_, err = remote.RPC.SetLocal(context.Background(), req)

	if err != nil {
		node.sugar.Error(err)
	}

	node.store.Del(key)
}

// Join join an existing network via another node
// This implements the pseudocode for in fig. 5 of the new paper
// for concurrent joins
func (node *Node) Join(remote scale.RemoteNode) {
	node.sugar.Infof("joining network via node at %s", remote.GetAddr())
	predecessor, err := remote.FindPredecessor(node.id)

	if err != nil {
		node.sugar.Fatal(err)
	} else if predecessor == nil {
		node.sugar.Fatal("no predecessor found")
	}

	s := NewRemoteNode(predecessor.GetAddr())
	p := s

	for !keyspace.BetweenRightInclusive(node.id, p.ID, s.ID) && !keyspace.Equal(p.ID, s.ID) {
		p = s
		successor, err := p.RPC.GetSuccessor(context.Background(), &pb.Empty{})

		if err != nil {
			node.sugar.Fatal(err)
		}

		s = NewRemoteNode(successor.GetAddr())
	}

	node.successor = s
	node.predecessor = p

	p.RPC.Notify(context.Background(), &pb.RemoteNode{Id: node.id[:], Addr: node.addr, Present: true})
	s.RPC.Notify(context.Background(), &pb.RemoteNode{Id: node.id[:], Addr: node.addr, Present: true})

	node.bootstrap(s)

	node.sugar.Info("joined network")
}

func (node *Node) bootstrap(n *RemoteNode) {
	for i := 0; i < scale.M; i++ {
		start := FingerMath(node.id[:], i, scale.M)
		p, _ := n.RPC.FindSuccessor(context.Background(), &pb.RemoteQuery{Id: start})
		var s *RemoteNode

		for bytes.Compare(p.GetId(), start) >= 0 {
			pred := NewRemoteNode(p.GetAddr())
			s = pred
			p, _ = pred.RPC.GetSuccessor(context.Background(), &pb.Empty{})
		}

		node.fingerTable[i] = s
	}
}

// GetLocal return a value stored on this node
func (node *Node) GetLocal(key scale.Key) ([]byte, error) {
	return node.store.Get(key), nil
}

// SetLocal set a value in the local store
func (node *Node) SetLocal(key scale.Key, value []byte) error {
	return node.store.Set(key, value)
}

// Get return a value stored on this node
func (node *Node) Get(key scale.Key) ([]byte, error) {
	succ, err := node.FindSuccessor(key)
	remoteNode := NewRemoteNode(succ.GetAddr())

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
	succ, err := node.FindSuccessor(key)
	remoteNode := NewRemoteNode(succ.GetAddr())

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
	node.mutex.RLock()
	defer node.mutex.RUnlock()

	var closestPrecedingRPC *pb.RemoteNode
	var err error

	if keyspace.Equal(node.id, node.predecessor.ID) {
		return NewRemoteNode(node.addr), nil
	}

	closestPreceding := NewRemoteNode(node.addr)
	predecessorSuccessor, _ := node.GetSuccessor()

	if keyspace.Equal(closestPreceding.ID, predecessorSuccessor.GetID()) {
		return closestPreceding, nil
	}

	for !keyspace.BetweenRightInclusive(id, closestPreceding.ID, predecessorSuccessor.GetID()) {
		closestPrecedingRPC, err = closestPreceding.RPC.ClosestPrecedingFinger(context.Background(), &pb.RemoteQuery{Id: id[:]})

		if err != nil {
			node.sugar.Fatal(err)
		}

		closestPreceding = NewRemoteNode(closestPrecedingRPC.Addr)
		predRes, _ := closestPreceding.RPC.GetSuccessor(context.Background(), &pb.Empty{})
		predecessorSuccessor = NewRemoteNode(predRes.GetAddr())
	}

	return closestPreceding, nil
}

// FindSuccessor returns the successor for this node
func (node *Node) FindSuccessor(id scale.Key) (scale.RemoteNode, error) {
	predecessor, err := node.FindPredecessor(id)

	if err != nil {
		node.sugar.Fatal(err)
	} else if predecessor.GetID() == node.GetID() {
		return node.successor, nil
	}

	predecessorRemode := NewRemoteNode(predecessor.GetAddr())
	successor, err := predecessorRemode.RPC.GetSuccessor(context.Background(), &pb.Empty{})

	if err != nil {
		node.sugar.Fatal(err)
	}

	return NewRemoteNode(successor.Addr), nil
}

// ClosestPrecedingFinger returns the closest preceding finger to the id
func (node *Node) ClosestPrecedingFinger(id scale.Key) (scale.RemoteNode, error) {
	node.mutex.RLock()
	defer node.mutex.RUnlock()

	for i := scale.M - 1; i >= 0; i-- {
		finger := node.fingerTable[i]

		if keyspace.Between(finger.ID, node.id, id) {
			return finger, nil
		}
	}
	return NewRemoteNode(node.addr), nil

}

// GetPredecessor returns the node's predecessor
func (node *Node) GetPredecessor() (scale.RemoteNode, error) {
	node.mutex.RLock()
	defer node.mutex.RUnlock()

	if node.predecessor != nil {
		return node.predecessor, nil
	}

	return nil, nil
}

// GetSuccessor retunrs the node's successor
func (node *Node) GetSuccessor() (scale.RemoteNode, error) {
	node.mutex.RLock()
	defer node.mutex.RUnlock()

	if node.successor != nil {
		return node.successor, nil
	}

	return nil, errors.New("no successor found")
}

// Shutdown leave the network
func (node *Node) Shutdown() {
	node.mutex.RLock()
	defer node.mutex.RUnlock()

	close(node.shutdownChannel)

	node.logger.Sync()

	if !keyspace.Equal(node.id, node.successor.ID) {
	}

	for _, remoteConnection := range node.remoteConnections {
		remoteConnection.clientConnection.Close()
	}
}

func (node *Node) stabilize() {
	if keyspace.Equal(node.successor.ID, node.id) {
		return
	}

	x, err := node.predecessor.RPC.GetSuccessor(context.Background(), &pb.Empty{})

	if err != nil {
		node.sugar.Fatal(err)
	}

	if x.Present && keyspace.Between(keyspace.ByteArrayToKey(x.Id), node.predecessor.ID, node.id) {
		node.mutex.Lock()
		node.predecessor = NewRemoteNode(x.Addr)
		node.mutex.Unlock()
	}

	x, err = node.successor.RPC.GetPredecessor(context.Background(), &pb.Empty{})

	if err != nil {
		node.sugar.Error(err)
	}

	if x.Present && keyspace.Between(keyspace.ByteArrayToKey(x.Id), node.id, node.successor.ID) {
		node.mutex.Lock()
		node.fingerTable[0] = NewRemoteNode(x.Addr)
		node.mutex.Unlock()
	}
}

func (node *Node) checkPredecessor() {
	node.mutex.Lock()
	defer node.mutex.Unlock()
	predecessor := node.predecessor
	id := predecessor.GetID()

	if predecessor == nil || keyspace.Equal(id, node.GetID()) {
		return
	}

	_, err := predecessor.RPC.Ping(context.Background(), &pb.Empty{})

	if err != nil {
		node.predecessor = nil
	}
}

// Notify is called when another node thinks it is our predecessor
func (node *Node) Notify(id scale.Key, addr string) error {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	if keyspace.Equal(node.id, node.successor.ID) && keyspace.Equal(node.id, node.predecessor.ID) {
		remote := NewRemoteNode(addr)
		node.fingerTable[0] = remote
		node.successor = remote
		node.predecessor = remote
		node.bootstrap(remote)

		return nil
	}

	if keyspace.Between(id, node.id, node.successor.ID) {
		successor := NewRemoteNode(addr)
		node.fingerTable[0] = successor
		node.successor = successor
		node.bootstrap(successor)
	}

	if keyspace.Between(id, node.predecessor.ID, node.id) {
		predecessor := NewRemoteNode(addr)
		node.predecessor = predecessor
		node.bootstrap(predecessor)
	}

	return nil
}

func (node *Node) fixNextFinger(next int) int {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	nextHash := FingerMath(node.id[:], next, scale.M)
	successor, _ := node.FindSuccessor(keyspace.ByteArrayToKey(nextHash))
	successorRemote := successor.(*RemoteNode)
	finger := NewRemoteNode(successorRemote.Addr)
	node.fingerTable[next] = finger
	return next + 1
}
