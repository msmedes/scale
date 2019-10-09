package scale

import pb "github.com/msmedes/scale/internal/app/scale/proto"

// RemoteNode contains metadata about another node
type RemoteNode struct {
	ID   Key
	Addr string
	RPC  pb.ScaleClient
}

// NewRemoteNode creates a new RemoteNode with an RPC client
func NewRemoteNode(addr string, node *Node) *RemoteNode {
	return &RemoteNode{ID: GenerateKey(addr), Addr: addr, RPC: GetScaleClient(addr, node)}
}
