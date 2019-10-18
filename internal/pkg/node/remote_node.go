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
