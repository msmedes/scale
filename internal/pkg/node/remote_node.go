package node

import (
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

// NewRemoteNode creates a new RemoteNode with an RPC client.
// This will reuse RPC connections if given the same address
func NewRemoteNode(addr string, node *Node) *RemoteNode {
	var remoteNode *RemoteNode
	var ok bool

	id := keyspace.GenerateKey(addr)
	remoteNode, ok = node.remoteConnections[id]

	if !ok {
		client, conn := rpc.NewClient(addr)
		remoteNode = &RemoteNode{ID: id, Addr: addr, RPC: client, clientConnection: conn}
		node.remoteConnections[id] = remoteNode
	}

	return remoteNode
}
