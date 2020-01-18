package node

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"sync"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/msmedes/scale/internal/pkg/keyspace"
	"github.com/msmedes/scale/internal/pkg/scale"
	"github.com/msmedes/scale/internal/pkg/store"
	"github.com/msmedes/scale/internal/pkg/trace"
	pb "github.com/msmedes/scale/internal/pkg/trace/proto"
	"go.uber.org/zap"
)

// StabilizeInterval sets how often to execute stabilization in ms
const StabilizeInterval = 1000

// FixFingerInterval sets how often to execute finger table updates in ms
const FixFingerInterval = 100

// Node main node class
type Node struct {
	scale.Node

	id              scale.Key
	addr            string
	port            string
	predecessor     scale.RemoteNode
	successor       scale.RemoteNode
	fingerTable     *FingerTable
	store           *store.MemoryStore
	logger          *zap.Logger
	sugar           *zap.SugaredLogger
	shutdownChannel chan struct{}
	mutex           sync.RWMutex
	traceClient     pb.TraceClient
	traceConn       *grpc.ClientConn
}

// NewNode create a new node
func NewNode(addr string) *Node {
	port := addr[strings.LastIndex(addr, ":")+1:]

	node := &Node{
		id:              keyspace.GenerateKey(addr),
		addr:            addr,
		port:            port,
		store:           store.NewMemoryStore(),
		shutdownChannel: make(chan struct{}),
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
	node.traceClient, node.traceConn = trace.NewClient("0.0.0.0:5000")

	node.SetupCloseHandler()
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

	for _, v := range node.fingerTable.Table {
		keys = append(keys, v.GetID())
	}

	return keys
}

// StabilizationStart runs periodically to ping the predecessor and
// stabilize the node.
func (node *Node) StabilizationStart() {
	ticker := time.NewTicker(StabilizeInterval * time.Millisecond)

	for {
		select {
		case <-ticker.C:
			node.checkPredecessor()
			node.stabilize()
		case <-node.shutdownChannel:
			ticker.Stop()
			return
		}
	}
}

// FixFingerStart runs periodically to make sure the finger table is up to date
func (node *Node) FixFingerStart() {
	next := 0
	ticker := time.NewTicker(FixFingerInterval * time.Millisecond)

	for {
		select {
		case <-ticker.C:
			next = node.fixNextFinger(next)
		case <-node.shutdownChannel:
			ticker.Stop()
			return
		}
	}
}

// TransferKeys transfer keys to the given node
func (node *Node) TransferKeys(id scale.Key, addr string) (count int) {
	count = 0
	remote := newRemoteNode(addr)

	if keyspace.Equal(id, node.id) {
		return
	}

	for _, k := range node.store.Keys() {
		if keyspace.GTE(k, id) {
			node.transferKey(k, remote)
			count++
		}
	}
	return
}

func (node *Node) transferKey(key scale.Key, remote scale.RemoteNode) {
	ctx := context.Background()
	val, err := node.GetLocal(ctx, key)

	if err != nil {
		node.sugar.Error(err)
	}

	err = remote.SetLocal(ctx, key, val)

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

func (node *Node) setSuccessor(successor scale.RemoteNode) {
	// fmt.Printf("setting successor to %s pred: %s\n", successor.GetAddr(), node.predecessor.GetAddr())
	node.successor = successor
}

func (node *Node) setPredecessor(predecessor scale.RemoteNode) {
	// fmt.Printf("setting predecessor to %s succ: %s\n", predecessor.GetAddr(), node.successor.GetAddr())
	node.predecessor = predecessor
}

func (node *Node) join(remote scale.RemoteNode) {
	node.mutex.RLock()
	defer node.mutex.RUnlock()

	ctx := context.Background()

	node.sugar.Infof("joining network via node at %s", remote.GetAddr())
	s, err := remote.FindPredecessor(ctx, node.id)

	if err != nil {
		node.sugar.Fatal(err)
	} else if s == nil {
		node.sugar.Fatal("no predecessor found")
	}

	p := s
	s, err = node.CallFunction(ctx, "GetSuccessor", p)

	if err != nil {
		node.sugar.Error(err)
	}

	for !keyspace.BetweenRightInclusive(node.id, p.GetID(), s.GetID()) && !keyspace.Equal(p.GetID(), s.GetID()) {
		p = s
		successor, err := p.GetSuccessor(ctx)

		if err != nil {
			node.sugar.Fatal(err)
		}
		s = successor
	}

	node.fingerTable.Table[0] = s
	node.setSuccessor(s)
	node.setPredecessor(p)

	p.Notify(node)
	s.Notify(node)
	node.bootstrap(s)

	node.sugar.Info("joined network")
}

func (node *Node) bootstrap(n scale.RemoteNode) {
	// var p, s scale.RemoteNode
	// var err error
	ctx := context.Background()

	for i := 0; i < scale.M; i++ {
		startKey := keyspace.ByteArrayToKey(node.fingerTable.starts[i])
		// // fmt.Printf("startKey %s\n", keyspace.KeyToString(startKey))
		// p, err = node.CallFunction(ctx, "FindSuccessor", n, startKey)

		// if err != nil {
		// 	node.sugar.Fatal(err)
		// } else if p == nil {
		// 	node.sugar.Fatalf("no successor found")
		// }
		// // startKey 00 1d 1c 11 where it breaks
		// s = p
		// for keyspace.GT(p.GetID(), startKey) {
		// 	s = p
		// 	// fmt.Printf("s %s\n", keyspace.KeyToString(s.GetID()))

		// 	p, err = node.CallFunction(ctx, "GetPredecessor", p)
		// 	// fmt.Printf("p %s\n", keyspace.KeyToString(p.GetID()))
		// 	if err != nil {
		// 		node.sugar.Fatal(err)
		// 	} else if p == nil {
		// 		node.sugar.Fatalf("no successor found")
		// 	}
		// }

		successor, _ := node.FindSuccessor(ctx, startKey)
		finger := newRemoteNode(successor.GetAddr())
		node.fingerTable.Table[i] = finger
	}
}

// GetLocal return a value stored on this node
func (node *Node) GetLocal(ctx context.Context, key scale.Key) ([]byte, error) {
	return node.store.Get(key), nil
}

// SetLocal set a value in the local store
func (node *Node) SetLocal(ctx context.Context, key scale.Key, value []byte) error {
	return node.store.Set(key, value)
}

// Get return a value stored on this node
func (node *Node) Get(ctx context.Context, key scale.Key) ([]byte, error) {
	var (
		val []byte
		err error
	)
	succ, err := node.FindSuccessor(ctx, key)
	if keyspace.Equal(succ.GetID(), node.id) {
		val, err = node.GetLocal(ctx, key)
	} else {
		val, err = succ.GetLocal(ctx, key)
	}

	if err != nil {
		return nil, err
	}

	return val, nil
}

// Set set a value in the local store
func (node *Node) Set(ctx context.Context, key scale.Key, value []byte) error {
	succ, err := node.FindSuccessor(ctx, key)
	if keyspace.Equal(node.id, succ.GetID()) {
		err = node.SetLocal(ctx, key, value)
	} else {
		remoteNode := newRemoteNode(succ.GetAddr())
		err = remoteNode.SetLocal(ctx, key, value)
	}

	if err != nil {
		return err
	}

	return nil
}

//FindPredecessor finds the predecessor to the id
func (node *Node) FindPredecessor(ctx context.Context, key scale.Key) (scale.RemoteNode, error) {
	node.mutex.RLock()
	defer node.mutex.RUnlock()

	var n1 scale.RemoteNode
	var err error

	n1 = node.toRemoteNode()
	successor := node.successor

	if !keyspace.BetweenRightInclusive(key, node.predecessor.GetID(), node.successor.GetID()) {
		n1, err = node.CallFunction(ctx, "ClosestPrecedingFinger", n1, key)
	}

	successor, err = node.CallFunction(ctx, "GetSuccessor", n1)
	if err != nil {
		return nil, err
	}

	successor = newRemoteNode(successor.GetAddr())

	if err != nil {
		return nil, err
	}

	for !keyspace.BetweenRightInclusive(key, n1.GetID(), successor.GetID()) && !keyspace.Equal(n1.GetID(), node.id) {
		n1, err = node.CallFunction(ctx, "ClosestPrecedingFinger", n1, key)
		if err != nil {
			return nil, err
		}

		successor, err = node.CallFunction(ctx, "GetSuccessor", n1)

		if err != nil {
			return nil, err
		}
	}

	return n1, nil
}

// FindSuccessor returns the successor for this node
func (node *Node) FindSuccessor(ctx context.Context, key scale.Key) (scale.RemoteNode, error) {
	predecessor, err := node.FindPredecessor(ctx, key)

	if err != nil {
		node.sugar.Fatal(err)
	} else if keyspace.Equal(predecessor.GetID(), node.GetID()) {
		return node.successor, nil
	}

	return predecessor.GetSuccessor(ctx)
}

// ClosestPrecedingFinger returns the closest preceding finger to the id
func (node *Node) ClosestPrecedingFinger(ctx context.Context, id scale.Key) (scale.RemoteNode, error) {
	node.mutex.RLock()
	defer node.mutex.RUnlock()

	for i := scale.M - 1; i >= 0; i-- {
		finger := node.fingerTable.Table[i]

		if keyspace.Between(finger.GetID(), node.id, id) {
			return finger, nil
		}
	}

	return node.toRemoteNode(), nil
}

// GetPredecessor returns the node's predecessor
func (node *Node) GetPredecessor(ctx context.Context) (scale.RemoteNode, error) {
	node.mutex.RLock()
	defer node.mutex.RUnlock()

	if node.predecessor != nil {
		return node.predecessor, nil
	}

	return nil, errors.New("no predecessor found")
}

// GetSuccessor retunrs the node's successor
func (node *Node) GetSuccessor(ctx context.Context) (scale.RemoteNode, error) {

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
		successorAddr := node.successor.GetAddr()
		predecessorAddr := node.predecessor.GetAddr()

		node.sugar.Infof("Transferring keys to %v", successorAddr)
		node.TransferKeys(node.predecessor.GetID(), successorAddr)

		node.sugar.Infof("Notifying %v of new predecessor %v", successorAddr, predecessorAddr)
		node.successor.SetPredecessor(predecessorAddr, node.GetAddr())

		node.sugar.Infof("Notifying %v of new sucessor %v", predecessorAddr, successorAddr)
		node.predecessor.SetSuccessor(successorAddr, node.GetAddr())
	}

	for _, remoteConnection := range remotes.data {
		node.sugar.Infof("Closing connection to %s", remoteConnection.GetAddr())
		remoteConnection.CloseConnection()
	}
	node.traceConn.Close()
	node.sugar.Info("Closing connection to trace server")
}

func (node *Node) stabilize() {
	var x scale.RemoteNode
	var err error

	ctx := context.Background()
	x, err = node.CallFunction(ctx, "GetSuccessor", node.predecessor)

	if err != nil {
		node.sugar.Fatal(err)
	}

	if x != nil && keyspace.Between(x.GetID(), node.predecessor.GetID(), node.id) {
		node.mutex.Lock()
		node.setPredecessor(newRemoteNode(x.GetAddr()))
		node.sugar.Infof("predecessor set to %v", node.predecessor.GetAddr())
		node.mutex.Unlock()
	}

	x, err = node.CallFunction(ctx, "GetPredecessor", node.successor)
	if err != nil {
		node.sugar.Error(err)
	}

	if x != nil && keyspace.Between(x.GetID(), node.id, node.successor.GetID()) {
		node.mutex.Lock()
		successor := newRemoteNode(x.GetAddr())
		node.fingerTable.Table[0] = successor
		node.setSuccessor(successor)
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
	}
}

// Notify is called when another node thinks it is our predecessor
func (node *Node) Notify(id scale.Key, addr string) error {

	if keyspace.Equal(node.id, node.successor.GetID()) && keyspace.Equal(node.id, node.predecessor.GetID()) {
		node.mutex.Lock()
		remote := newRemoteNode(addr)
		node.fingerTable.Table[0] = remote
		node.setSuccessor(remote)
		node.setPredecessor(remote)
		node.sugar.Infof("Predecessor and successor set to %s", remote.GetAddr())

		node.sugar.Infof("Transferring keys to %s", remote.GetAddr())
		numKeysTransferred := node.TransferKeys(remote.GetID(), remote.GetAddr())
		node.sugar.Infof("Transferred %v keys", numKeysTransferred)

		node.mutex.Unlock()
		node.bootstrap(remote)

		return nil
	}

	if keyspace.Between(id, node.id, node.successor.GetID()) {
		node.mutex.Lock()
		successor := newRemoteNode(addr)
		node.fingerTable.Table[0] = successor
		node.setSuccessor(successor)
		node.sugar.Infof("Successor set to %v %p", node.successor.GetID(), node.successor)

		node.sugar.Infof("transferring keys to %s", successor.GetAddr())
		count := node.TransferKeys(successor.GetID(), successor.GetAddr())
		node.sugar.Infof("Transferred %v keys", count)

		node.mutex.Unlock()
		node.bootstrap(successor)
	}

	if keyspace.Between(id, node.predecessor.GetID(), node.id) {
		node.mutex.Lock()
		predecessor := newRemoteNode(addr)
		node.setPredecessor(predecessor)
		node.sugar.Infof("Predecessor set to %+v %p", node.predecessor, node.predecessor)
		node.mutex.Unlock()
		node.bootstrap(predecessor)
	}

	return nil
}

//ToRemoteNode convert to remote node
func (node *Node) toRemoteNode() scale.RemoteNode {
	return newRemoteNodeWithID(node.addr, node.id)
}

// Check figure 5 for init_finger_table, something is up
func (node *Node) fixNextFinger(next int) int {
	node.mutex.RLock()
	defer node.mutex.RUnlock()

	// fmt.Printf("next %d\r", next)

	ctx := context.Background()

	if next == scale.M {
		next = 0
		// fmt.Printf("pred: %s %s succ: %s %s\n", node.predecessor.GetAddr(), keyspace.KeyToString(node.predecessor.GetID()), node.successor.GetAddr(), keyspace.KeyToString(node.successor.GetID()))
		// for i := range node.fingerTable.Table {
		// 	keyStart := keyspace.KeyToString(keyspace.ByteArrayToKey(node.fingerTable.starts[i]))
		// 	finger := keyspace.KeyToString(node.fingerTable.Table[i].GetID())
		// 	fmt.Printf("%d start: %s, finger: %s\n", i, keyStart, finger)
		// }
	}

	nextHash := keyspace.ByteArrayToKey(node.fingerTable.starts[next])
	successor, _ := node.FindSuccessor(ctx, nextHash)
	finger := newRemoteNode(successor.GetAddr())
	node.fingerTable.Table[next] = finger

	return next + 1
}

// GetKeys returns all the keys in the store
func (node *Node) GetKeys() []string {
	return node.store.KeysAsString()
}

// CallFunction is a handy little function to remove all the
// if keyspace.Equal(node.id, remote.ID) boilerplate
func (node *Node) CallFunction(ctx context.Context, funcName string, remote scale.RemoteNode, params ...interface{}) (scale.RemoteNode, error) {
	var (
		in          []reflect.Value
		nodeReflect reflect.Value
		function    reflect.Value
		response    []reflect.Value
		err         error
	)

	// Call was rejecting the length of the returned []reflect.Value for
	// functions with no params when these few lines were in their own
	// function, but it works this way for some reason?
	if len(params) == 0 {
		in = []reflect.Value{reflect.ValueOf(ctx)}
	} else {
		in = make([]reflect.Value, len(params))
		for k, param := range params {
			in[k] = reflect.ValueOf(param)
		}
		in = prependContext(ctx, in)
	}

	if keyspace.Equal(node.GetID(), remote.GetID()) {
		nodeReflect = reflect.ValueOf(node)
		function, err = functionFactory(nodeReflect, funcName, len(in))
	} else {
		nodeReflect = reflect.ValueOf(remote)
		function, err = functionFactory(nodeReflect, funcName, len(in))
	}
	if err != nil {
		return nil, err
	}

	response = function.Call(in)
	remoteNode := response[0].Interface()

	if response[1].Interface() != nil {
		responseErr := response[1].Interface().(error)
		return nil, responseErr
	}

	return remoteNode.(scale.RemoteNode), nil
}

// lol what is this java
func functionFactory(value reflect.Value, functionName string, numParams int) (reflect.Value, error) {
	function := value.MethodByName(fmt.Sprintf("%s", functionName))
	numIn := function.Type().NumIn()
	if numParams != numIn {
		return reflect.ValueOf(nil), fmt.Errorf("invalid number of arguments. got: %v, want: %v", numParams, numIn)
	}
	return function, nil
}

// SetupCloseHandler creates a listener to notify the node if it receives an
// interrupt signal from the OS and run the shutdown prodecure.
func (node *Node) SetupCloseHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		node.sugar.Info("Keyboard interrupt, shutting down")
		node.Shutdown()
		node.sugar.Info("K bye!")
		os.Exit(0)
	}()
}

// SetSuccessor sets the successor
func (node *Node) SetSuccessor(succAddr string, clientAddr string) error {
	node.mutex.RLock()
	defer node.mutex.RUnlock()
	// fmt.Printf("successor set to %s\n", succAddr)
	node.successor = newRemoteNode(succAddr)
	remotes.data[clientAddr].CloseConnection()
	node.fingerTable.Table[0] = node.successor
	node.bootstrap(node.successor)
	return nil
}

// SetPredecessor sets the predecessor
func (node *Node) SetPredecessor(predAddr string, clientAddr string) error {
	node.mutex.RLock()
	defer node.mutex.RUnlock()
	// fmt.Printf("predecessor set to %s \n", predAddr)
	node.setPredecessor(newRemoteNode(predAddr))
	remotes.data[clientAddr].CloseConnection()
	return nil
}

// don't @ me
func prependContext(ctx context.Context, in []reflect.Value) []reflect.Value {
	in = append(in, reflect.ValueOf(0))
	copy(in[1:], in)
	in[0] = reflect.ValueOf(ctx)
	return in
}

// SendTraceID extracts the traceID from incoming context and sends
// a message to the trace server that it has received the traceID
func (node *Node) SendTraceID(ctx context.Context, functionCall string) context.Context {
	md, _ := metadata.FromIncomingContext(ctx)
	if traceID, ok := md["traceid"]; ok {
		_, err := node.traceClient.AppendTrace(context.Background(), &pb.AppendTraceRequest{
			TraceID:      traceID[0],
			Addr:         node.GetAddr(),
			FunctionCall: functionCall,
			Timestamp:    time.Now().UnixNano(),
		})
		if err != nil {
			node.sugar.Info(err)
		}
	}
	return metadata.NewIncomingContext(ctx, md)
}
