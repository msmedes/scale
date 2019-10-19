package node

import (
	"context"

	"github.com/msmedes/scale/internal/pkg/keyspace"
	"github.com/msmedes/scale/internal/pkg/rpc"
	pb "github.com/msmedes/scale/internal/pkg/rpc/proto"
	"github.com/msmedes/scale/internal/pkg/scale"
	"google.golang.org/grpc"
)

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
func (r *RemoteNode) FindPredecessor(key scale.Key) (scale.RemoteNode, error) {
	predecessor, err := r.RPC.FindPredecessor(
		context.Background(),
		&pb.RemoteQuery{Id: key[:]},
	)

	if err != nil {
		return nil, err
	}

	return NewRemoteNode(predecessor.GetAddr()), nil
}

// FindSuccessor proxy
func (r *RemoteNode) FindSuccessor(key scale.Key) (scale.RemoteNode, error) {
	p, err := r.RPC.FindSuccessor(context.Background(), &pb.RemoteQuery{Id: key[:]})

	if err != nil {
		return nil, err
	}

	return NewRemoteNode(p.GetAddr()), nil
}

// NewRemoteNode creates a new RemoteNode with an RPC client.
// This will reuse RPC connections if given the same address
func NewRemoteNode(addr string) *RemoteNode {
	id := keyspace.GenerateKey(addr)
	client, conn := rpc.NewClient(addr)

	return &RemoteNode{
		ID:               id,
		Addr:             addr,
		RPC:              client,
		clientConnection: conn,
	}
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
func (r *RemoteNode) GetSuccessor() (scale.RemoteNode, error) {
	successor, err := r.RPC.GetSuccessor(context.Background(), &pb.Empty{})

	if err != nil {
		return nil, err
	}

	return NewRemoteNode(successor.GetAddr()), nil
}

//GetPredecessor proxy
func (r *RemoteNode) GetPredecessor() (scale.RemoteNode, error) {
	predecessor, err := r.RPC.GetPredecessor(context.Background(), &pb.Empty{})

	if err != nil {
		return nil, err
	}

	return NewRemoteNode(predecessor.GetAddr()), nil
}

//Ping proxy
func (r *RemoteNode) Ping() error {
	_, err := r.RPC.Ping(context.Background(), &pb.Empty{})

	return err
}
