package node

import (
	"context"
	"fmt"
	"os"
	"sync"

	"go.uber.org/zap"

	"github.com/msmedes/scale/internal/pkg/keyspace"
	"github.com/msmedes/scale/internal/pkg/rpc"
	pb "github.com/msmedes/scale/internal/pkg/rpc/proto"
	"github.com/msmedes/scale/internal/pkg/scale"
	"google.golang.org/grpc"
)

type remotesCache struct {
	mutex sync.RWMutex
	data  map[string]scale.RemoteNode
	sugar *zap.SugaredLogger
}

func (r *remotesCache) get(addr string) scale.RemoteNode {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.data[addr]
}

func (r *remotesCache) set(addr string, node scale.RemoteNode) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.data[addr] = node
}

func newLogger() *zap.SugaredLogger {
	logger, _ := zap.NewDevelopment(
		zap.Fields(
			zap.String("remotes", fmt.Sprintf("pid: %d", os.Getpid())),
		),
	)
	sugar := logger.Sugar()
	return sugar
}

var remotes *remotesCache = &remotesCache{
	data:  make(map[string]scale.RemoteNode),
	sugar: newLogger(),
}

// RemoteNode contains metadata about another node
type RemoteNode struct {
	scale.RemoteNode

	ID               scale.Key
	Addr             string
	RPC              pb.ScaleClient
	clientConnection *grpc.ClientConn
}

// GetID getter for ID
func (r *RemoteNode) GetID() scale.Key {
	return r.ID
}

// GetAddr getter for address
func (r *RemoteNode) GetAddr() string {
	return r.Addr
}

// FindPredecessor proxy for RPC call
func (r *RemoteNode) FindPredecessor(ctx context.Context, key scale.Key) (scale.RemoteNode, error) {
	predecessor, err := r.RPC.FindPredecessor(
		ctx,
		&pb.RemoteQuery{Id: key[:]},
	)

	if err != nil {
		return nil, err
	}

	return newRemoteNode(predecessor.GetAddr()), nil
}

// FindSuccessor proxy
func (r *RemoteNode) FindSuccessor(ctx context.Context, key scale.Key) (scale.RemoteNode, error) {
	p, err := r.RPC.FindSuccessor(ctx, &pb.RemoteQuery{Id: key[:]})

	if err != nil {
		return nil, err
	}

	return newRemoteNode(p.GetAddr()), nil
}

// NewRemoteNode creates a new RemoteNode with an RPC client.
// This will reuse RPC connections if given the same address
func newRemoteNode(addr string) scale.RemoteNode {
	id := keyspace.GenerateKey(addr)
	return newRemoteNodeWithID(addr, id)
}

func newRemoteNodeWithID(addr string, id scale.Key) scale.RemoteNode {
	var remote scale.RemoteNode
	remote = remotes.get(addr)

	if remote != nil {
		return remote
	}

	client, conn := rpc.NewClient(addr)

	remote = &RemoteNode{
		ID:               id,
		Addr:             addr,
		RPC:              client,
		clientConnection: conn,
	}

	remotes.set(addr, remote)

	return remote
}

// CloseConnection closes the client connection
func (r *RemoteNode) CloseConnection() error {
	err := r.clientConnection.Close()
	if err != nil {
		return err
	}
	delete(remotes.data, r.GetAddr()) // I guess this can be handled here?
	return nil
}

// Notify proxy
func (r *RemoteNode) Notify(node scale.Node) error {
	id := node.GetID()
	_, err := r.RPC.Notify(
		context.Background(),
		&pb.RemoteNode{Id: id[:], Addr: node.GetAddr(), Present: true},
	)

	return err
}

//GetSuccessor proxy
func (r *RemoteNode) GetSuccessor(ctx context.Context) (scale.RemoteNode, error) {
	successor, err := r.RPC.GetSuccessor(ctx, &pb.Empty{})

	if err != nil {
		return nil, err
	}

	return newRemoteNode(successor.GetAddr()), nil
}

//GetPredecessor proxy
func (r *RemoteNode) GetPredecessor(ctx context.Context) (scale.RemoteNode, error) {
	predecessor, err := r.RPC.GetPredecessor(ctx, &pb.Empty{})

	if err != nil {
		return nil, err
	}

	return newRemoteNode(predecessor.GetAddr()), nil
}

//Ping proxy
func (r *RemoteNode) Ping() error {
	_, err := r.RPC.Ping(context.Background(), &pb.Empty{})

	return err
}

//GetLocal proxy
func (r *RemoteNode) GetLocal(ctx context.Context, key scale.Key) ([]byte, error) {
	val, err := r.RPC.GetLocal(
		ctx,
		&pb.GetRequest{Key: key[:]},
	)

	if err != nil {
		return nil, err
	}

	return val.GetValue(), nil
}

//SetLocal proxy
func (r *RemoteNode) SetLocal(ctx context.Context, key scale.Key, val []byte) error {
	_, err := r.RPC.SetLocal(
		ctx,
		&pb.SetRequest{Key: key[:], Value: val},
	)

	return err
}

//ClosestPrecedingFinger proxy
func (r *RemoteNode) ClosestPrecedingFinger(ctx context.Context, id scale.Key) (scale.RemoteNode, error) {
	res, err := r.RPC.ClosestPrecedingFinger(ctx, &pb.RemoteQuery{Id: id[:]})

	if err != nil {
		return nil, err
	}

	return newRemoteNode(res.GetAddr()), nil
}

// SetPredecessor proxy
func (r *RemoteNode) SetPredecessor(addr string, clientAddr string) error {
	_, err := r.RPC.SetPredecessor(context.Background(), &pb.ShutdownRequest{
		EssorAddr:  addr,
		ClientAddr: clientAddr,
	})

	return err
}

// SetSuccessor proxy
func (r *RemoteNode) SetSuccessor(addr string, clientAddr string) error {
	_, err := r.RPC.SetSuccessor(context.Background(), &pb.ShutdownRequest{
		EssorAddr:  addr,
		ClientAddr: clientAddr,
	})

	return err
}

// GetNetwork proxy
func (r *RemoteNode) GetNetwork(nodes []string) ([]string, error) {
	network, err := r.RPC.GetNetwork(context.Background(), &pb.NetworkMessage{
		Nodes: nodes,
	})

	return network.Nodes, err
}

func clearRemotes() {
	remotes.mutex.Lock()
	defer remotes.mutex.Unlock()

	for k := range remotes.data {
		remote := remotes.data[k]

		err := remote.CloseConnection()
		if err != nil {
			remotes.sugar.Error(err)
		}
		delete(remotes.data, k)
	}
}
