package node

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/msmedes/scale/internal/pkg/keyspace"
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
	predecessor       scale.RemoteNode
	successor         scale.RemoteNode
	fingerTable       table
	store             *store.MemoryStore
	logger            *zap.Logger
	sugar             *zap.SugaredLogger
	remoteConnections map[scale.Key]*RemoteNode
	shutdownChannel   chan struct{}
	mutex             sync.RWMutex
	functions         map[string]reflect.Value
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
			node.fixNextFinger(next)
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
	s, err := remote.FindPredecessor(node.id)

	if err != nil {
		node.sugar.Fatal(err)
	} else if s == nil {
		node.sugar.Fatal("no predecessor found")
	}

	p := s
	if keyspace.Equal(p.GetID(), node.id) {
		s, err = node.GetSuccessor()
	} else {
		s, err = p.GetSuccessor()
	}
	if err != nil {
		node.sugar.Fatal("no successor to predecessor found")
	}

	// ok this is what the new function call would look like
	// the last arg are the params, we can just wrap them in interfaces
	// like make([]reflect.Value{}{{key}, {node}}) or just make a helper function
	// to create the interface

	// remoteTest, err := node.CallCommand(p, "GetSuccessor", make([]reflect.Value, 0))
	// if err != nil {
	// 	node.sugar.Info("nope")
	// }
	// node.sugar.Infof("holy shit did that work? %+v", remoteTest) // it did

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
		startKey := keyspace.ByteArrayToKey(start)

		if keyspace.Equal(n.GetID(), node.id) {
			p, err = node.FindSuccessor(startKey)
		} else {
			p, err = n.FindSuccessor(startKey)
		}

		if err != nil {
			node.sugar.Fatal(err)
		} else if p == nil {
			node.sugar.Fatalf("no successor found")
		}

		s := p

		for keyspace.GT(p.GetID(), startKey) {
			s = p

			if keyspace.Equal(p.GetID(), node.id) {
				p, err = node.GetPredecessor()
			} else {
				p, err = p.GetPredecessor()
			}

			if err != nil {
				node.sugar.Fatal(err)
			} else if p == nil {
				node.sugar.Fatalf("no successor found")
			}
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
	var (
		val []byte
		err error
	)
	succ, err := node.FindSuccessor(key)
	if keyspace.Equal(succ.GetID(), node.id) {
		val, err = node.GetLocal(key)
	} else {
		val, err = succ.GetLocal(key)
	}

	if err != nil {
		return nil, err
	}

	return val, nil
}

// Set set a value in the local store
func (node *Node) Set(key scale.Key, value []byte) error {
	succ, err := node.FindSuccessor(key)
	if keyspace.Equal(node.id, succ.GetID()) {
		err = node.SetLocal(key, value)
	} else {
		remoteNode := newRemoteNode(succ.GetAddr())
		err = remoteNode.SetLocal(key, value)
	}

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
		if keyspace.Equal(n1.GetID(), node.id) {
			n1, err = node.ClosestPrecedingFinger(key)
		} else {
			n1, err = n1.ClosestPrecedingFinger(key)
		}
	}

	if keyspace.Equal(n1.GetID(), node.id) {
		successor, err = node.GetSuccessor()
	} else {
		successor, err = n1.GetSuccessor()
	}
	if err != nil {
		return nil, err
	}

	successor = newRemoteNode(successor.GetAddr())

	if err != nil {
		return nil, err
	}

	for !keyspace.BetweenRightInclusive(key, n1.GetID(), successor.GetID()) && !keyspace.Equal(n1.GetID(), node.id) {
		if keyspace.Equal(n1.GetID(), node.id) {
			n1, err = node.ClosestPrecedingFinger(key)
		} else {
			n1, err = n1.ClosestPrecedingFinger(key)
		}
		if err != nil {
			return nil, err
		}

		if keyspace.Equal(n1.GetID(), node.id) {
			successor, err = node.GetSuccessor()
		} else {
			successor, err = n1.GetSuccessor()
		}

		if err != nil {
			return nil, err
		}
	}

	return n1, nil
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
	var x scale.RemoteNode
	var err error

	if keyspace.Equal(node.id, node.predecessor.GetID()) {
		x, err = node.GetSuccessor()
	} else {
		x, err = node.predecessor.GetSuccessor()
	}

	if err != nil {
		node.sugar.Fatal(err)
	}

	if x != nil && keyspace.Between(x.GetID(), node.predecessor.GetID(), node.id) {
		node.mutex.Lock()
		node.predecessor = newRemoteNode(x.GetAddr())
		node.sugar.Infof("predecessor set to %v", node.predecessor.GetAddr())
		node.mutex.Unlock()
	}

	if keyspace.Equal(node.successor.GetID(), node.id) {
		x, err = node.GetSuccessor()
	} else {
		x, err = node.successor.GetSuccessor()
	}

	if err != nil {
		node.sugar.Error(err)
	}

	if x != nil && keyspace.Between(x.GetID(), node.id, node.successor.GetID()) {
		node.mutex.Lock()
		successor := newRemoteNode(x.GetAddr())
		node.fingerTable[0] = successor
		node.successor = successor
		node.mutex.Unlock()
		node.sugar.Infof("successor set to %v", successor.GetAddr())
	}
}

func (node *Node) checkPredecessor() {
	node.mutex.RLock()
	defer node.mutex.RUnlock()

	predecessor := node.predecessor
	id := predecessor.GetID()

	if predecessor == nil || keyspace.Equal(id, node.GetID()) {
		return
	}

	err := predecessor.Ping()

	if err != nil {
		node.predecessor = nil
		predecessor.CloseConnection()
		delete(remotes.data, predecessor.GetAddr())
	}
}

// Notify is called when another node thinks it is our predecessor
func (node *Node) Notify(id scale.Key, addr string) error {

	if keyspace.Equal(node.id, node.successor.GetID()) && keyspace.Equal(node.id, node.predecessor.GetID()) {
		node.mutex.Lock()
		remote := newRemoteNode(addr)
		node.fingerTable[0] = remote
		node.successor = remote
		node.predecessor = remote
		node.sugar.Infof("predecessor and successor set to %+v %p", remote, remote)
		node.mutex.Unlock()
		node.bootstrap(remote)

		return nil
	}

	if keyspace.Between(id, node.id, node.successor.GetID()) {
		node.mutex.Lock()
		successor := newRemoteNode(addr)
		node.fingerTable[0] = successor
		node.successor = successor
		node.sugar.Infof("successor set to %+v %p", node.successor, node.successor)
		node.mutex.Unlock()
		node.bootstrap(successor)
	}

	if keyspace.Between(id, node.predecessor.GetID(), node.id) {
		node.mutex.Lock()
		predecessor := newRemoteNode(addr)
		node.predecessor = predecessor
		node.sugar.Infof("predecessor set to %+v %p", node.predecessor, node.predecessor)
		node.mutex.Unlock()
		node.bootstrap(predecessor)
	}

	return nil
}

//ToRemoteNode convert to remote node
func (node *Node) toRemoteNode() scale.RemoteNode {
	return newRemoteNodeWithID(node.addr, node.id)
}

func (node *Node) fixNextFinger(next int) int {
	node.mutex.RLock()
	defer node.mutex.RUnlock()

	nextHash := fingerMath(node.id[:], next)
	successor, _ := node.FindSuccessor(keyspace.ByteArrayToKey(nextHash))
	finger := newRemoteNode(successor.GetAddr())
	node.fingerTable[next] = finger
	return next + 1
}

// GetKeys returns all the keys in the store
func (node *Node) GetKeys() []string {
	return node.store.KeysAsString()
}

// CallCommand should theoretically call the correct function on the correct
// Node or RemoteNode object aka Mike's first foray into metaprogramming
func (node *Node) CallCommand(remote scale.RemoteNode, funcName string, in []reflect.Value) (*RemoteNode, error) {
	var (
		nodeReflect reflect.Value
		function    reflect.Value
		response    []reflect.Value
		err         error
	)

	if keyspace.Equal(node.GetID(), remote.GetID()) {
		nodeReflect = reflect.ValueOf(node)
		node.sugar.Infof("%+v", nodeReflect)
		function, err = functionFactory(nodeReflect, funcName, len(in))
		if err != nil {
			return nil, err
		}
	} else {
		nodeReflect = reflect.ValueOf(remote)
		function, err = functionFactory(nodeReflect, funcName, len(in))
		if err != nil {
			return nil, err
		}
	}
	// Ok we should unpack the type in the function that calls it, I'm thinking
	// but there's probably a way to call TypeOf(reponse[0]) or something in the cast
	// so we can actually use this function in testing
	// Oh wait, we could return an interface with whatever in it and then cast on the other end

	response = function.Call(in)
	remoteNode := response[0].Interface().(*RemoteNode) //maybe call .(TypeOf(response[0]) here?)

	if response[1].Interface() != nil {
		responseErr := response[1].Interface().(error)
		return nil, responseErr
	}

	return remoteNode, nil

}

// lol what is this java
func functionFactory(value reflect.Value, functionName string, numParams int) (reflect.Value, error) {
	function := value.MethodByName(fmt.Sprintf("%s", functionName))
	t := function.Type()
	if numParams != t.NumIn() {
		return reflect.ValueOf(nil), fmt.Errorf("invalid number of arguments. got: %v, want: %v", numParams, t.NumIn())
	}
	return function, nil
}
