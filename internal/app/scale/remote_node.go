package scale

import pb "github.com/msmedes/scale/internal/app/scale/proto"

// RemoteNode contains metadata about another node
type RemoteNode struct {
	ID   Key
	Addr string
	RPC  pb.ScaleClient
}

// NewRemoteNode creates a new RemoteNode with an RPC client.
// This will reuse RPC connections if given the same address
func NewRemoteNode(addr string, node *Node) *RemoteNode {
	id := GenerateKey(addr)
	remoteNode, ok := node.remoteConnections[id]

	if !ok {
		return &RemoteNode{ID: id, Addr: addr, RPC: GetScaleClient(addr, node)}
	}

	return remoteNode
}
