package node

import (
	"strings"

	"github.com/msmedes/scale/internal/pkg/keyspace"
	"github.com/msmedes/scale/internal/pkg/rpc"
	pb "github.com/msmedes/scale/internal/pkg/rpc/proto"
	"github.com/msmedes/scale/internal/pkg/scale"
)

// RemoteNode contains metadata about another node
type RemoteNode struct {
	scale.RemoteNode

	ID   scale.Key
	Addr string
	RPC  pb.ScaleClient
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
	port := addr[strings.LastIndex(addr, ":")+1:]
	id := keyspace.GenerateKey(port)
	remoteNode, ok := node.remoteConnections[id]

	if !ok {
		return &RemoteNode{ID: id, Addr: addr, RPC: rpc.NewClient(addr)}
	}

	return remoteNode
}
