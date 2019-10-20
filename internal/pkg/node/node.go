package node

import (
	"bytes"
	"errors"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/msmedes/scale/internal/pkg/keyspace"
	"github.com/msmedes/scale/internal/pkg/scale"
	"github.com/msmedes/scale/internal/pkg/store"
	"go.uber.org/zap"
)

// StabilizeInterval how often to execute stabilization in seconds
const StabilizeInterval = 5

// Node main node class
type Node struct {
	scale.Node

	id                scale.Key
	addr              string
	port              string
	predecessor       scale.RemoteNode
	successor         scale.RemoteNode
	fingerTable       table
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

	node.fingerTable = newFingerTable(node.toRemoteNode())
	node.successor = node.toRemoteNode()
	node.predecessor = node.toRemoteNode()

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
		keys = append(keys, k.GetID())
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
	remote := newRemoteNode(addr)

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

func (node *Node) transferKey(key scale.Key, remote scale.RemoteNode) {
	val, err := node.GetLocal(key)

	if err != nil {
		node.sugar.Error(err)
	}

	err = remote.SetLocal(key, val)

	if err != nil {
		node.sugar.Error(err)
		return
	}

	node.store.Del(key)
}

// JoinAddr join an existing network via another node
// This implements the pseudocode for in fig. 5 of the new paper
// for concurrent joins
func (node *Node) JoinAddr(addr string) {
	node.join(newRemoteNode(addr))
}

func (node *Node) join(remote scale.RemoteNode) {
	node.sugar.Infof("joining network via node at %s", remote.GetAddr())
	predecessor, err := remote.FindPredecessor(node.id)

	if err != nil {
		node.sugar.Fatal(err)
	} else if predecessor == nil {
		node.sugar.Fatal("no predecessor found")
	}

	s := predecessor
	p := s

	for !keyspace.BetweenRightInclusive(node.id, p.GetID(), s.GetID()) && !keyspace.Equal(p.GetID(), s.GetID()) {
		p = s
		successor, err := p.GetSuccessor()

		if err != nil {
			node.sugar.Fatal(err)
		}

		s = successor
	}

	node.successor = s
	node.predecessor = p

	p.Notify(node)
	s.Notify(node)

	node.bootstrap(s)

	node.sugar.Info("joined network")
}

func (node *Node) bootstrap(n scale.RemoteNode) {
	var p scale.RemoteNode
	var err error

	for i := 0; i < scale.M; i++ {
		start := fingerMath(node.id[:], i)

		if keyspace.Equal(n.GetID(), node.id) {
			node.sugar.Info("node bootstrap")
			p, err = node.FindSuccessor(keyspace.ByteArrayToKey(start))
		} else {
			node.sugar.Info("n bootstrap")
			p, err = n.FindSuccessor(keyspace.ByteArrayToKey(start))
		}

		if err != nil {
			node.sugar.Fatal(err)
		} else if p == nil {
			node.sugar.Fatalf("no successor found")
		}

		s := p
		pID := p.GetID()

		for bytes.Compare(pID[:], start[:]) > 0 {
			pred := p
			s = pred
			p, err = pred.GetSuccessor()

			if err != nil {
				node.sugar.Fatal(err)
			} else if p == nil {
				node.sugar.Fatalf("no successor found")
			}

			pID = p.GetID()
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
	remoteNode := succ
	val, err := remoteNode.GetLocal(key)

	if err != nil {
		return nil, err
	}

	return val, nil
}

// Set set a value in the local store
func (node *Node) Set(key scale.Key, value []byte) error {
	succ, err := node.FindSuccessor(key)
	remoteNode := newRemoteNode(succ.GetAddr())
	err = remoteNode.SetLocal(key, value)

	if err != nil {
		return err
	}

	return nil
}

//FindPredecessor finds the predecessor to the id
func (node *Node) FindPredecessor(key scale.Key) (scale.RemoteNode, error) {
	node.mutex.RLock()
	defer node.mutex.RUnlock()

	var n1 scale.RemoteNode
	var err error

	n1 = node.toRemoteNode()
	successor := node.successor

	if !keyspace.BetweenRightInclusive(key, node.predecessor.GetID(), node.successor.GetID()) {
		n1, err = node.ClosestPrecedingFinger(key)
		if err != nil {
			return nil, err
		}
	} else {
		return n1, nil
	}

	// ok let's get the successor of n1

	successor, err = n1.GetSuccessor()
	if err != nil {
		return nil, err
	}

	for !keyspace.BetweenRightInclusive(key, n1.GetID(), successor.GetID()) && !keyspace.Equal(n1.GetID(), node.id) {
		n1, err = n1.ClosestPrecedingFinger(key)
		if err != nil {
			return nil, err
		}

		successor, err = n1.GetSuccessor()

		if err != nil {
			return nil, err
		}
	}

	return newRemoteNode(n1.GetAddr()), nil
}

// FindSuccessor returns the successor for this node
func (node *Node) FindSuccessor(key scale.Key) (scale.RemoteNode, error) {
	predecessor, err := node.FindPredecessor(key)

	if err != nil {
		node.sugar.Fatal(err)
	} else if keyspace.Equal(predecessor.GetID(), node.GetID()) {
		return node.successor, nil
	}

	successor, err := predecessor.GetSuccessor()

	if err != nil {
		node.sugar.Fatal(err)
	}

	return successor, nil
}

// ClosestPrecedingFinger returns the closest preceding finger to the id
func (node *Node) ClosestPrecedingFinger(id scale.Key) (scale.RemoteNode, error) {
	node.mutex.RLock()
	defer node.mutex.RUnlock()

	for i := scale.M - 1; i >= 0; i-- {
		finger := node.fingerTable[i]

		if keyspace.Between(finger.GetID(), node.id, id) {
			return finger, nil
		}
	}
	return node.toRemoteNode(), nil
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

	node.logger.Sync()
	node.sugar.Sync()

	close(node.shutdownChannel)

	if !keyspace.Equal(node.id, node.successor.GetID()) {
	}

	for _, remoteConnection := range node.remoteConnections {
		remoteConnection.clientConnection.Close()
	}
}

func (node *Node) stabilize() {
	if keyspace.Equal(node.successor.GetID(), node.id) {
		node.sugar.Info("stabilize: node and successor are same")
	}

	x, err := node.predecessor.GetSuccessor()

	if err != nil {
		node.sugar.Fatal(err)
	}

	if x != nil && keyspace.Between(x.GetID(), node.predecessor.GetID(), node.id) {
		node.mutex.Lock()
		node.predecessor = newRemoteNode(x.GetAddr())
		node.sugar.Infof("predecessor set to %+v", node.predecessor)
		node.mutex.Unlock()
	}

	x, err = node.successor.GetPredecessor()

	if err != nil {
		node.sugar.Error(err)
	}

	if x != nil && keyspace.Between(x.GetID(), node.id, node.successor.GetID()) {
		node.mutex.Lock()
		node.fingerTable[0] = newRemoteNode(x.GetAddr())
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

	err := predecessor.Ping()

	if err != nil {
		node.predecessor = nil
	}
}

// Notify is called when another node thinks it is our predecessor
func (node *Node) Notify(id scale.Key, addr string) error {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	if keyspace.Equal(node.id, node.successor.GetID()) && keyspace.Equal(node.id, node.predecessor.GetID()) {
		remote := newRemoteNode(addr)
		node.fingerTable[0] = remote
		node.successor = remote
		node.predecessor = remote
		node.sugar.Infof("predecessor and successor setto %+v", remote)
		node.bootstrap(remote)

		return nil
	}

	if keyspace.Between(id, node.id, node.successor.GetID()) {
		successor := newRemoteNode(addr)
		node.fingerTable[0] = successor
		node.successor = successor
		node.sugar.Infof("successor set to %+v", node.successor)
		node.bootstrap(successor)
	}

	if keyspace.Between(id, node.predecessor.GetID(), node.id) {
		predecessor := newRemoteNode(addr)
		node.predecessor = predecessor
		node.sugar.Infof("predecessor set to %+v", node.predecessor)
		node.bootstrap(predecessor)
	}

	return nil
}

//ToRemoteNode convert to remote node
func (node *Node) toRemoteNode() scale.RemoteNode {
	return newRemoteNodeWithID(node.addr, node.id)
}

func (node *Node) fixNextFinger(next int) int {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	nextHash := fingerMath(node.id[:], next)
	successor, _ := node.FindSuccessor(keyspace.ByteArrayToKey(nextHash))
	successorRemote := successor.(*RemoteNode)
	finger := newRemoteNode(successorRemote.Addr)
	node.fingerTable[next] = finger
	return next + 1
}
